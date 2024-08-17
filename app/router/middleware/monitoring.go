package middleware

import (
	"fmt"
	"net/http"

	"github.com/RyoheiTomiyama/phraze-api/domain/infra/monitoring"
	"github.com/RyoheiTomiyama/phraze-api/util/env"
)

func Monitoring(conf *env.Config, mHttp monitoring.IHttp) func(http.Handler) http.Handler {
	if err := mHttp.Setup(&monitoring.Options{
		Dsn: conf.Sentry.DSN,
	}); err != nil {
		fmt.Println(err)
	}

	return mHttp.Handler
}
