package monitoring

import (
	"context"
	"net/http"
)

type Options struct {
	Dsn string
}

type IClient interface {
	ReportError(ctx context.Context, err error)
	// RecordEvent(ctx context.Context)
}

type IHttp interface {
	Setup(opt *Options) error
	Handler(next http.Handler) http.Handler
}
