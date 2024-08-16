package monitoring

import (
	"context"
	"fmt"

	"github.com/RyoheiTomiyama/phraze-api/domain/infra/monitoring"
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

func (c *client) RecordEvent(ctx context.Context, l monitoring.Level, m string, arg ...any) {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		data := map[string]interface{}{}
		for i := 0; i < len(arg); i += 2 {
			key := fmt.Sprintf("%v", arg[i])
			value := fmt.Sprintf("%v", arg[i+1])
			data[key] = value
		}

		hub.AddBreadcrumb(&sentry.Breadcrumb{
			Type:     string(l),
			Category: string(l),
			Message:  m,
			Level:    toSentryLevel(l),
			Data:     data,
		}, nil)
	}
}

func toSentryLevel(l monitoring.Level) sentry.Level {
	switch l {
	case monitoring.LevelInfo:
		return sentry.LevelInfo
	case monitoring.LevelWarning:
		return sentry.LevelWarning
	case monitoring.LevelError:
		return sentry.LevelError
	default:
		return sentry.LevelDebug
	}
}
