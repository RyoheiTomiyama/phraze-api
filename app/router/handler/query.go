package handler

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/RyoheiTomiyama/phraze-api/router/graph/generated"
	"github.com/RyoheiTomiyama/phraze-api/router/graph/resolver"
)

func PostQuery() http.HandlerFunc {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{}}))
	return srv.ServeHTTP
}
