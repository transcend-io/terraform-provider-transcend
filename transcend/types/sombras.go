package types

import (
	graphql "github.com/hasura/go-graphql-client"
)

type Sombra struct {
	ID  graphql.String `json:"id"`
	URL graphql.String `json:"url"`
}
