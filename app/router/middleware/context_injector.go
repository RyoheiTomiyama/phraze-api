package middleware

import (
	"net/http"

	"github.com/RyoheiTomiyama/phraze-api/util/env"
	"github.com/RyoheiTomiyama/phraze-api/util/logger"
)

// Bearerから認証情報をcontextに詰めるMiddleware
func ContextInjector(conf *env.Config, l logger.ILogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = conf.WithCtx(ctx)
			ctx = l.WithCtx(ctx)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
