package logger

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"

	"github.com/RyoheiTomiyama/phraze-api/infra/monitoring"
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
	Debug(msg string, arg ...any)
	Info(msg string, arg ...any)
	Warning(msg string, arg ...any)
	Error(err error, arg ...any)
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

func (l *logger) Debug(msg string, arg ...any) {
	l.logger.Debug(msg, arg...)
}

func (l *logger) Info(msg string, arg ...any) {
	l.logger.Info(msg, arg...)
}

func (l *logger) Warning(msg string, arg ...any) {
	l.logger.Warn(msg, arg...)
}

func (l *logger) error(err error, arg ...any) {
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
	l.logger.Error(err.Error(), arg...)
}
func (l *logger) Error(err error, arg ...any) {
	l.error(err, arg...)
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
	l.error(err, arg...)
	l.reportError(ctx, err)
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
