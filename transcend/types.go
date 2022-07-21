package transcend

import "github.com/shurcooL/graphql"

type DataPoint struct {
	Name  graphql.String
	Title struct {
		DefaultMessage graphql.String
	}
	Description struct {
		DefaultMessage graphql.String
	}
	DataCollection struct {
		VisualID graphql.String
	}
}

type DbIntegrationQuerySuggestionInput string

type RequestActionObjectResolver string

type DataPointSubDataPointInput struct {
	Name       graphql.String            `json:"name"`
	Desciption graphql.String            `json:"description"`
	Categories []DataSubCategoryInput    `json:"categories"`
	Purposes   []PurposeSubCategoryInput `json:"purposes"`
	Attributes []AttributeInput          `json:"attributes"`
}

type DataCategoryType string

type DataSubCategoryInput struct {
	Name     graphql.String   `json:"name"`
	Category DataCategoryType `json:"category"`
}

type ProcessingPurpose string

type PurposeSubCategoryInput struct {
	Name    graphql.String    `json:"name"`
	Purpose ProcessingPurpose `json:"purpose"`
}

type AttributeInput struct {
	Key    graphql.String   `json:"key"`
	Values []graphql.String `json:"values"`
}

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

type Header struct {
	Name     graphql.String  `json:"name"`
	Value    graphql.String  `json:"value"`
	IsSecret graphql.Boolean `json:"isSecret"`
}

type CustomHeaderInput Header

type DataSilo struct {
	ID      graphql.String `json:"id"`
	Title   graphql.String `json:"title"`
	Link    graphql.String `json:"link"`
	Type    graphql.String `json:"type"`
	Catalog struct {
		HasAvcFunctionality graphql.Boolean `json:"hasAvcFunctionality"`
	} `json:"catalog"`
}

type APIKey struct {
	ID     graphql.String
	Title  graphql.String
	Scopes []struct {
		Type graphql.String
	}
	DataSilos []DataSilo
}

type ScopeName string
