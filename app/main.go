package main

import (
	"net/http"
	"time"

	"github.com/RyoheiTomiyama/phraze-api/application/usecase/auth"
	"github.com/RyoheiTomiyama/phraze-api/application/usecase/card"
	"github.com/RyoheiTomiyama/phraze-api/application/usecase/deck"
	"github.com/RyoheiTomiyama/phraze-api/infra/db"
	firebaseAuth "github.com/RyoheiTomiyama/phraze-api/infra/firebase/auth"
	"github.com/RyoheiTomiyama/phraze-api/infra/gemini"
	"github.com/RyoheiTomiyama/phraze-api/infra/monitoring"
	"github.com/RyoheiTomiyama/phraze-api/router"
	"github.com/RyoheiTomiyama/phraze-api/router/graph/directive"
	"github.com/RyoheiTomiyama/phraze-api/router/graph/resolver"
	card_service "github.com/RyoheiTomiyama/phraze-api/service/card"
	"github.com/RyoheiTomiyama/phraze-api/util/env"
	"github.com/RyoheiTomiyama/phraze-api/util/logger"
)

func main() {
	l := logger.New(logger.Options{Level: logger.LevelDebug, Debug: true})

	config, err := env.New()
	if err != nil {
		panic(err)
	}

	// infra
	dbClient, err := db.NewClient(db.DataSourceOption{
		Host:     config.DB.HOST,
		Port:     config.DB.PORT,
		DBName:   config.DB.DB_NAME,
		User:     config.DB.USER,
		Password: config.DB.PASSWORD,
	})
	if err != nil {
		panic(err)
	}

	geminiClient, err := gemini.New(gemini.ClientOption{APIKey: config.Gemini.API_KEY})
	if err != nil {
		panic(err)
	}

	firebaseAuthClient, err := firebaseAuth.New()
	if err != nil {
		panic(err)
	}

	monitoringHttp := monitoring.NewHttp()
	monitoringClient := monitoring.New()
	l = l.WithMonitoring(monitoringClient)

	// service
	cardService := card_service.NewService()

	// usecase
	authUsecase := auth.New(firebaseAuthClient)
	cardUsecase := card.New(dbClient, geminiClient, cardService)
	deckUsecase := deck.New(dbClient)

	resolver := resolver.New(cardUsecase, deckUsecase)
	directive := directive.New()

	r := router.New(config, resolver, &directive, l, monitoringHttp, authUsecase)

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
