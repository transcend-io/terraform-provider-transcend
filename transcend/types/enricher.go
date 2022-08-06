package types

import graphql "github.com/hasura/go-graphql-client"

type Enricher struct {
	ID              graphql.String
	Title           graphql.String
	Description     graphql.String
	Url             graphql.String
	InputIdentifier struct {
		ID graphql.String
	}
	Identifiers []struct {
		ID graphql.String
	}
	Actions []RequestAction
	Headers []Header
}
