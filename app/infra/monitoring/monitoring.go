package monitoring

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/domain/infra/monitoring"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/getsentry/sentry-go"
)

type client struct {
}

func New() monitoring.IClient {
	return &client{}
}

func (c *client) ReportError(ctx context.Context, err error) {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.CaptureException((err))
	}
}

func (c *client) Setup(opt *monitoring.Options) error {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: opt.Dsn,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 0.1,
	}); err != nil {
		return errutil.Wrap(err)
	}

	return nil
}
