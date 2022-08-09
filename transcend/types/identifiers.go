package types

import (
	graphql "github.com/hasura/go-graphql-client"
)

type Identifier struct {
	ID   graphql.String `json:"id"`
	Name graphql.String `json:"name"`
}

// func listIdentifiers() []Identifier {
// 	client := NewClient("https://api.dev.trancsend.com/graphql", os.Getenv("TRANSCEND_KEY"))

// 	var query struct {
// 		Identifiers []Identifier `graphql:"identifiers(first: 3)"`
// 	}

// 	err := client.graphql.Query(context.Background(), &query, map[string]interface{}{}, graphql.OperationName("ApiKey"))
// 	assert.Nil(t, err)

// 	return query.Enricher
// }
