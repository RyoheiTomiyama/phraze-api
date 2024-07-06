package router

import (
	"github.com/RyoheiTomiyama/phraze-api/application/usecase/auth"
	"github.com/RyoheiTomiyama/phraze-api/router/handler"
	"github.com/RyoheiTomiyama/phraze-api/router/middleware"
	"github.com/RyoheiTomiyama/phraze-api/util/env"
	"github.com/RyoheiTomiyama/phraze-api/util/logger"
	"github.com/go-chi/chi/v5"
)

type router struct {
	config      *env.Config
	logger      logger.ILogger
	authUsecase auth.IAuthUsecase
}

type IRouter interface {
	Handler() *chi.Mux
}

func New(config *env.Config, l logger.ILogger, authUsecase auth.IAuthUsecase) IRouter {
	return &router{config, l, authUsecase}
}

func (r *router) Handler() *chi.Mux {
	chiRouter := chi.NewRouter()

	chiRouter.Use(
		middleware.ContextInjector(r.config, r.logger),
		middleware.Authrization(r.authUsecase),
	)
	// chiRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("welcome phraze"))
	// })

	chiRouter.Get("/playground", handler.Playground())
	chiRouter.Post("/query", handler.PostQuery())

	return chiRouter
}
