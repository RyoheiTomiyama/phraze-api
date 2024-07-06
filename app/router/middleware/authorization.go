package middleware

import (
	"net/http"

	"github.com/RyoheiTomiyama/phraze-api/application/usecase/auth"
	authutil "github.com/RyoheiTomiyama/phraze-api/util/auth"
	"github.com/RyoheiTomiyama/phraze-api/util/logger"
)

// Bearerから認証情報をcontextに詰めるMiddleware
func Authrization(authUsecase auth.IAuthUsecase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			token := r.Header.Get("Authorization")
			if token == "" {
				next.ServeHTTP(w, r)

				return
			}

			user, err := authUsecase.ParseToken(ctx, token)
			if err != nil {
				l := logger.FromCtx(ctx)
				l.Error(err, "token", token)
				http.Error(w, err.Error(), http.StatusForbidden)

				return
			}

			ctx = authutil.New(user).WithCtx(ctx)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
