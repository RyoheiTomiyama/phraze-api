package logger

import (
	"context"
	"log/slog"
	"os"
)

type Level = slog.Level

const (
	LevelDebug Level = -4
	LevelInfo  Level = 0
	LevelWarn  Level = 4
	LevelError Level = 8
)

type logger struct {
	logger *slog.Logger
	// Notify NotifyFunc
}

type Options struct {
	Level Level
}

func New(opt Options) *logger {
	return &logger{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: opt.Level,
		})),
	}
}

type loggerCtxKey struct{}

// 引数で与えられたロガーを context に詰め、新たな context を返す
func (l *logger) WithCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerCtxKey{}, l)
}

// context からロガーを取り出す。取り出せない場合はデフォルトのロガーを返す
func FromCtx(ctx context.Context) *logger {
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

func (l *logger) Error(msg string, arg ...any) {
	l.logger.Error(msg, arg...)
}
