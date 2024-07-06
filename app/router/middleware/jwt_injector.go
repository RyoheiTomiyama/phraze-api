package middleware

import (
	"net/http"

	"github.com/RyoheiTomiyama/phraze-api/util/logger"
)

func JWTInjector(jwtHandler jwt.HandlerInterface) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				// JWT を必要としないエンドポイントも存在しているため、エラーにはしない
				next.ServeHTTP(w, r)

				return
			}

			ctx := r.Context()
			payload, err := jwtHandler.Parse(token)
			if err != nil {
				l := logger.FromCtx(ctx)
				l.Error("jwt の payload を取得するのに失敗しました", "error", err.Error())
				http.Error(w, err.Error(), http.StatusForbidden)

				return
			}

			r = r.WithContext(jwtHandler.WithCtx(ctx, *payload))

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
