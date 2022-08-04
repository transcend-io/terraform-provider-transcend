package types

import "github.com/shurcooL/graphql"

type APIKey struct {
	ID     graphql.String `json:"id"`
	Title  graphql.String `json:"title"`
	Scopes []struct {
		Type graphql.String `json:"type"`
	} `json:"scopes"`
}
