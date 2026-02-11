package gql

import (
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// NewHandler создает HTTP handler для GraphQL с GraphiQL
func NewHandler(schema *graphql.Schema) http.Handler {
	return handler.New(&handler.Config{
		Schema:   schema,
		GraphiQL: true,
		Pretty:   true,
	})
}
