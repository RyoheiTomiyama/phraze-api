package middleware

import (
	"net/http"

	"github.com/RyoheiTomiyama/phraze-api/util/env"
	"github.com/go-chi/cors"
)

func CorsHandler(conf *env.Config) func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins: conf.APP.CORS,
		AllowedHeaders: []string{"Origin", "X-Requested-With", "Content-Type", "Authorization", "Accept"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
	})
}
