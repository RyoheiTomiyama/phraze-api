package monitoring

import (
	"context"

	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/getsentry/sentry-go"
)

type client struct {
}

type IClient interface {
	ReportError(ctx context.Context, err error)
}

func New() IClient {
	return &client{}
}

func (c *client) ReportError(ctx context.Context, err error) {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.CaptureException((err))
	}
}

type MonitoringOptions struct {
	Dsn string
}

func Setup(opt *MonitoringOptions) error {
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
