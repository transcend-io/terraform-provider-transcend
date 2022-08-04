package types

import "github.com/shurcooL/graphql"

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
