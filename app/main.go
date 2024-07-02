package main

import (
	"context"
	"net/http"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/infra/db"
	"github.com/RyoheiTomiyama/phraze-api/infra/router"
	"github.com/RyoheiTomiyama/phraze-api/util/env"
	"github.com/RyoheiTomiyama/phraze-api/util/logger"
)

func main() {
	ctx := context.Background()

	l := logger.New(logger.Options{Level: logger.LevelDebug, Debug: true})
	ctx = l.WithCtx(ctx)

	config, err := env.New(ctx)
	if err != nil {
		panic(err)
	}
	ctx = config.WithCtx(ctx)

	db, err := db.NewClient(config.DB.DSN)
	if err != nil {
		panic(err)
	}

	deck, err := db.GetDeck(ctx, 1)
	if err != nil {
		panic(err)
	}

	l.Debug("deck", "d", &deck)

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
