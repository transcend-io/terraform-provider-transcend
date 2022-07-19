package transcend

import "github.com/shurcooL/graphql"

// Planning to move all the types here

type RequestAction string

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
