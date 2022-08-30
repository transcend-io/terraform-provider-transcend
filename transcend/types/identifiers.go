package types

import (
	graphql "github.com/hasura/go-graphql-client"
)

type Identifier struct {
	ID   graphql.String `json:"id"`
	Name graphql.String `json:"name"`
}
