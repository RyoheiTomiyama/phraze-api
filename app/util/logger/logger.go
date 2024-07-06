package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
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
	logger    *slog.Logger
	debugMode bool
	// Notify NotifyFunc
}

type ILogger interface {
	WithCtx(ctx context.Context) context.Context
	Debug(msg string, arg ...any)
	Info(msg string, arg ...any)
	Warning(msg string, arg ...any)
	Error(err error, arg ...any)
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

func (l *logger) Error(err error, arg ...any) {
	if l.debugMode {
		st := errutil.ErrorWithStackTrace(err)
		if len(st) > 0 {
			arg = append(arg, "stack_trace", st)
		}
	}
	l.logger.Error(err.Error(), arg...)
}
