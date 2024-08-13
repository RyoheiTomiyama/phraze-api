package monitoring

import (
	"net/http"

	"github.com/RyoheiTomiyama/phraze-api/util/auth"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/samber/lo"
)

func Handler(next http.Handler) http.Handler {
	sentryHandler := sentryhttp.New(sentryhttp.Options{})
	return sentryHandler.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := auth.FromCtx(ctx)

		if hub := sentry.GetHubFromContext(ctx); hub != nil {
			if user != nil {
				hub.Scope().SetUser(sentry.User{
					ID:       user.ID,
					Username: lo.FromPtrOr(user.Name, "no name"),
				})
			}
		}
		next.ServeHTTP(w, r)
	})
}
