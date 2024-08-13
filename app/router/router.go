package router

import (
	"github.com/RyoheiTomiyama/phraze-api/application/usecase/auth"
	"github.com/RyoheiTomiyama/phraze-api/router/graph/generated"
	"github.com/RyoheiTomiyama/phraze-api/router/graph/resolver"
	"github.com/RyoheiTomiyama/phraze-api/router/handler"
	"github.com/RyoheiTomiyama/phraze-api/router/middleware"
	"github.com/RyoheiTomiyama/phraze-api/util/env"
	"github.com/RyoheiTomiyama/phraze-api/util/logger"
	"github.com/go-chi/chi/v5"
)

type router struct {
	config      *env.Config
	resolver    *resolver.Resolver
	directive   *generated.DirectiveRoot
	logger      logger.ILogger
	authUsecase auth.IAuthUsecase
}

type IRouter interface {
	Handler() *chi.Mux
}

func New(
	config *env.Config,
	resolver *resolver.Resolver,
	directive *generated.DirectiveRoot,
	l logger.ILogger,
	authUsecase auth.IAuthUsecase,
) IRouter {
	return &router{config, resolver, directive, l, authUsecase}
}

func (r *router) Handler() *chi.Mux {
	chiRouter := chi.NewRouter()

	chiRouter.Use(
		middleware.Recoverer,
		middleware.CorsHandler(r.config),
		middleware.ContextInjector(r.config, r.logger),
		middleware.Authrization(r.authUsecase),
		middleware.Monitoring(r.config),
	)
	// chiRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("welcome phraze"))
	// })

	chiRouter.Get("/playground", handler.Playground())
	chiRouter.Post("/query", handler.PostQuery(r.resolver, r.directive))

	return chiRouter
}
