package types

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

type DataPointSubDataPointInput struct {
	Name       graphql.String            `json:"name"`
	Desciption graphql.String            `json:"description"`
	Categories []DataSubCategoryInput    `json:"categories"`
	Purposes   []PurposeSubCategoryInput `json:"purposes"`
	Attributes []AttributeInput          `json:"attributes"`
}

type PurposeSubCategoryInput struct {
	Name    graphql.String    `json:"name"`
	Purpose ProcessingPurpose `json:"purpose"`
}

type DataSubCategoryInput struct {
	Name     graphql.String   `json:"name"`
	Category DataCategoryType `json:"category"`
}

type AttributeInput struct {
	Key    graphql.String   `json:"key"`
	Values []graphql.String `json:"values"`
}
