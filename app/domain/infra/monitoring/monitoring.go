package monitoring

import (
	"context"
	"net/http"
)

type Level string

const (
	LevelDebug   Level = "debug"
	LevelInfo    Level = "info"
	LevelWarning Level = "warning"
	LevelError   Level = "error"
	j
)

type Options struct {
	Dsn string
}

type IClient interface {
	ReportError(ctx context.Context, err error)
	RecordEvent(ctx context.Context, l Level, m string, arg ...any)
}

type IHttp interface {
	Setup(opt *Options) error
	Handler(next http.Handler) http.Handler
}
