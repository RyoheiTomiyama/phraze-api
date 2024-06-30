package main

import (
	"context"
	"net/http"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/infra/router"
	"github.com/RyoheiTomiyama/phraze-api/util/logger"
)

func main() {
	ctx := context.Background()

	l := logger.New(logger.Options{Level: logger.LevelDebug, Debug: true})
	l.WithCtx(ctx)

	r := router.New()

	server := &http.Server{
		Addr:              ":3000",
		ReadHeaderTimeout: 1 * time.Second,
		ReadTimeout:       3 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       1 * time.Second,
		Handler:           r.Handler(),
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
