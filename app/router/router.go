package router

import (
	"github.com/RyoheiTomiyama/phraze-api/router/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type router struct{}

type IRouter interface {
	Handler() *chi.Mux
}

func New() IRouter {
	return &router{}
}

func (r *router) Handler() *chi.Mux {
	chiRouter := chi.NewRouter()

	chiRouter.Use(middleware.Logger)
	// chiRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("welcome phraze"))
	// })

	chiRouter.Get("/playground", handler.Playground())
	chiRouter.Post("/query", handler.PostQuery())

	return chiRouter
}
