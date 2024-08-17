package monitoring

import (
	"net/http"

	"github.com/RyoheiTomiyama/phraze-api/domain/infra/monitoring"
	"github.com/RyoheiTomiyama/phraze-api/util/auth"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/samber/lo"
)

type mhttp struct{}

func NewHttp() monitoring.IHttp {
	return &mhttp{}
}

func (m *mhttp) Setup(opt *monitoring.Options) error {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: opt.Dsn,
		// Debug: true,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1,
		// EnableTracing:    true,
	}); err != nil {
		return errutil.Wrap(err)
	}

	return nil
}

func (m *mhttp) Handler(next http.Handler) http.Handler {
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
