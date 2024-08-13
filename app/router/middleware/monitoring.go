package middleware

import (
	"fmt"
	"net/http"

	"github.com/RyoheiTomiyama/phraze-api/infra/monitoring"
	"github.com/RyoheiTomiyama/phraze-api/util/env"
)

func Monitoring(conf *env.Config) func(http.Handler) http.Handler {
	if err := monitoring.Setup(&monitoring.MonitoringOptions{
		Dsn: conf.Sentry.DSN,
	}); err != nil {
		fmt.Println(err)
	}

	return monitoring.Handler
}
