package handler

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/RyoheiTomiyama/phraze-api/infra/graph/generated"
	"github.com/RyoheiTomiyama/phraze-api/infra/graph/resolver"
)

func PostQuery() http.HandlerFunc {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{}}))
	return srv.ServeHTTP
}
