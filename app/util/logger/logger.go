package logger

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"

	"github.com/RyoheiTomiyama/phraze-api/domain/infra/monitoring"
	"github.com/golang-cz/devslog"
	"github.com/samber/lo"
)

type Level = slog.Level

const (
	LevelDebug Level = -4
	LevelInfo  Level = 0
	LevelWarn  Level = 4
	LevelError Level = 8
)

type logger struct {
	logger     *slog.Logger
	debugMode  bool
	monitoring monitoring.IClient
}

type ILogger interface {
	WithCtx(ctx context.Context) context.Context
	WithMonitoring(monitoring monitoring.IClient) ILogger

	Debug(ctx context.Context, msg string, arg ...any)
	Info(ctx context.Context, msg string, arg ...any)
	Warning(ctx context.Context, msg string, arg ...any)
	Error(ctx context.Context, err error, arg ...any)

	// レスポンスは正常終了させるけど、エラー通知したいときに使う
	ErrorWithNotify(ctx context.Context, err error, arg ...any)
}

type Options struct {
	Level Level
	Debug bool
}

func New(opt Options) ILogger {
	handler := lo.Ternary[slog.Handler](opt.Debug,
		devslog.NewHandler(os.Stdout, &devslog.Options{
			HandlerOptions: &slog.HandlerOptions{
				Level: opt.Level,
			},
		}),
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: opt.Level,
		}),
	)

	return &logger{
		logger:    slog.New(handler),
		debugMode: opt.Debug,
	}
}

type loggerCtxKey struct{}

// 引数で与えられたロガーを context に詰め、新たな context を返す
func (l *logger) WithCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerCtxKey{}, l)
}

func (l *logger) WithMonitoring(monitoring monitoring.IClient) ILogger {
	l.monitoring = monitoring
	return l
}

// context からロガーを取り出す。取り出せない場合はデフォルトのロガーを返す
func FromCtx(ctx context.Context) ILogger {
	l, ok := ctx.Value(loggerCtxKey{}).(*logger)
	if !ok {
		return New(Options{Level: LevelDebug})
	}

	return l
}

func (l *logger) Debug(ctx context.Context, msg string, arg ...any) {
	l.logger.DebugContext(ctx, msg, arg...)

	if l.monitoring != nil {
		l.monitoring.RecordEvent(ctx, monitoring.LevelDebug, msg, arg...)
	}
}

func (l *logger) Info(ctx context.Context, msg string, arg ...any) {
	l.logger.InfoContext(ctx, msg, arg...)

	if l.monitoring != nil {
		l.monitoring.RecordEvent(ctx, monitoring.LevelInfo, msg, arg...)
	}
}

func (l *logger) Warning(ctx context.Context, msg string, arg ...any) {
	l.logger.WarnContext(ctx, msg, arg...)

	if l.monitoring != nil {
		l.monitoring.RecordEvent(ctx, monitoring.LevelWarning, msg, arg...)
	}
}

func (l *logger) error(ctx context.Context, err error, arg ...any) {
	if l.debugMode {
		_, name, line, ok := runtime.Caller(2)
		if ok {
			arg = append(arg, "call", fmt.Sprintf("%s:%d", name, line))
		}

		st := l.createStackTrace(err)
		if len(st) > 0 {
			arg = append(arg, "stack_trace", st)
		}
	}
	l.logger.ErrorContext(ctx, err.Error(), arg...)
}
func (l *logger) Error(ctx context.Context, err error, arg ...any) {
	l.error(ctx, err, arg...)

	if l.monitoring != nil {
		l.monitoring.RecordEvent(ctx, monitoring.LevelError, err.Error(), arg...)
	}
}

func (l *logger) createStackTrace(err error) string {
	type stackError interface {
		StackTrace() []uintptr
	}

	var se stackError
	if errors.As(err, &se) {
		// log用にStackTraceを整形する
		frames := extractFrames(se.StackTrace(), 4)
		traceString := ""
		for _, f := range frames {
			file := f.File
			if strings.HasPrefix(file, "/go/src/app/") {
				relativeIndex := strings.Index(file, "/go/src/app/")
				file = file[relativeIndex+8:]
			}
			function := f.Function
			if i := strings.LastIndex(function, "/"); i > 0 {
				function = function[i+1:]
			}
			traceString += fmt.Sprintf("- %s\n  %s:%d\n", function, file, f.Line)
		}
		return traceString
	}

	return ""
}

func (l *logger) ErrorWithNotify(ctx context.Context, err error, arg ...any) {
	l.error(ctx, err, arg...)

	if l.monitoring != nil {
		l.reportError(ctx, err)
	}
}

func (l *logger) reportError(ctx context.Context, err error) {
	l.monitoring.ReportError(ctx, err)
}

// log吐き出し用にFrameに変換する
func extractFrames(pcs []uintptr, depth int) []runtime.Frame {
	var frames = make([]runtime.Frame, 0, len(pcs))
	callersFrames := runtime.CallersFrames(pcs[:min(len(pcs), depth)])

	for {
		callerFrame, more := callersFrames.Next()
		frames = append(frames, callerFrame)

		if !more {
			break
		}
	}

	return frames
}
