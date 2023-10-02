package graphql_api

import (
	"context"
)

func (g *GraphQLAPIClient) SetContext(ctx context.Context) {
	g.defaultContext = ctx
}

func (g *GraphQLAPIClient) Query(query interface{}, variables map[string]interface{}) error {
	return g.gqlClient.Query(g.defaultContext, query, variables)
}

func (g *GraphQLAPIClient) QueryRaw(query string, variables map[string]interface{}) ([]byte, error) {
	return g.gqlClient.QueryRaw(g.defaultContext, query, variables)
}

func (g *GraphQLAPIClient) Mutate(query string, variables map[string]interface{}) error {
	return g.gqlClient.Mutate(g.defaultContext, query, variables, nil)
}

func (g *GraphQLAPIClient) Exec(query string, variables map[string]interface{}) error {
	return g.gqlClient.Exec(g.defaultContext, query, variables, nil)
}

func (g *GraphQLAPIClient) ExecRaw(query string, variables map[string]interface{}) ([]byte, error) {
	return g.gqlClient.ExecRaw(g.defaultContext, query, variables)
}
