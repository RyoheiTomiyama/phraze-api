package handler

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
)

func Playground() http.HandlerFunc {
	return playground.Handler("GraphQL playground", "/query")
}
