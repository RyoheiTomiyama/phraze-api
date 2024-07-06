package main

import (
	"context"
	"net/http"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/infra/db"
	"github.com/RyoheiTomiyama/phraze-api/router"
	"github.com/RyoheiTomiyama/phraze-api/util/env"
	"github.com/RyoheiTomiyama/phraze-api/util/logger"
)

func main() {
	ctx := context.Background()

	l := logger.New(logger.Options{Level: logger.LevelDebug, Debug: true})
	ctx = l.WithCtx(ctx)

	config, err := env.New()
	if err != nil {
		panic(err)
	}
	ctx = config.WithCtx(ctx)

	_, err = db.NewClient(db.DataSourceOption{
		Host:     config.DB.HOST,
		Port:     config.DB.PORT,
		DBName:   config.DB.DB_NAME,
		User:     config.DB.USER,
		Password: config.DB.PASSWORD,
	})
	if err != nil {
		panic(err)
	}

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
