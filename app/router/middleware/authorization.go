package middleware

import (
	"net/http"
	"strings"

	"github.com/RyoheiTomiyama/phraze-api/application/usecase/auth"
	authutil "github.com/RyoheiTomiyama/phraze-api/util/auth"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/RyoheiTomiyama/phraze-api/util/logger"
)

// Bearerから認証情報をcontextに詰めるMiddleware
func Authrization(authUsecase auth.IAuthUsecase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.FromCtx(ctx)

			token := r.Header.Get("Authorization")
			if token == "" {
				next.ServeHTTP(w, r)

				return
			}

			// トークンの `Bearer` と `ey...` を分離する
			splitToken, ok := strings.CutPrefix(token, "Bearer ")
			if !ok {
				err := errutil.New(errutil.CodeForbidden, "bearer token の形式が不正です")
				http.Error(w, err.Error(), http.StatusForbidden)

				return
			}

			user, err := authUsecase.ParseToken(ctx, splitToken)
			if err != nil {
				l := logger.FromCtx(ctx)
				l.Error(err, "token", token)
				http.Error(w, err.Error(), http.StatusForbidden)

				return
			}

			log.Debug("Authorization", "user", user)

			ctx = authutil.New(user).WithCtx(ctx)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
