package types

import "github.com/shurcooL/graphql"

// Enums
type DbIntegrationQuerySuggestionInput string
type RequestActionObjectResolver string
type DataCategoryType string
type ProcessingPurpose string
type RequestAction string
type DataSiloConnectionState string
type ScopeName string

type Header struct {
	Name     graphql.String  `json:"name"`
	Value    graphql.String  `json:"value"`
	IsSecret graphql.Boolean `json:"isSecret"`
}

type CustomHeaderInput Header
