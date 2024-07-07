package middleware

import (
	"net/http"

	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/RyoheiTomiyama/phraze-api/util/logger"
)

func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				ctx := r.Context()
				log := logger.FromCtx(ctx)

				err := errutil.New(errutil.CodeInternalError, "panic!!: %+v", v)
				log.Error(err)

				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
