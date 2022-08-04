package types

import (
	"github.com/shurcooL/graphql"
)

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

func ToDataPointSubDataPointInputList(origs []interface{}) []DataPointSubDataPointInput {
	vals := make([]DataPointSubDataPointInput, len(origs))
	for i, orig := range origs {
		newVal := orig.(map[string]interface{})
		vals[i] = DataPointSubDataPointInput{
			graphql.String(newVal["name"].(string)),
			graphql.String(newVal["description"].(string)),
			ToDataSubCategoryInputList(newVal["categories"].([]interface{})),
			ToPurposeSubCategoryInputList(newVal["purposes"].([]interface{})),
			ToAttributeInputList(newVal["attributes"].([]interface{})),
		}
	}

	return vals
}

func ToDataSubCategoryInputList(origs []interface{}) []DataSubCategoryInput {
	vals := make([]DataSubCategoryInput, len(origs))
	for i, orig := range origs {
		newVal := orig.(map[string]interface{})
		vals[i] = DataSubCategoryInput{
			graphql.String(newVal["name"].(string)),
			DataCategoryType(newVal["category"].(string)),
		}
	}

	return vals
}

func ToPurposeSubCategoryInputList(origs []interface{}) []PurposeSubCategoryInput {
	vals := make([]PurposeSubCategoryInput, len(origs))
	for i, orig := range origs {
		newVal := orig.(map[string]interface{})
		vals[i] = PurposeSubCategoryInput{
			graphql.String(newVal["name"].(string)),
			ProcessingPurpose(newVal["purpose"].(string)),
		}
	}

	return vals
}

func ToAttributeInputList(origs []interface{}) []AttributeInput {
	vals := make([]AttributeInput, len(origs))
	for i, orig := range origs {
		newVal := orig.(map[string]interface{})
		vals[i] = AttributeInput{
			graphql.String(newVal["key"].(string)),
			ToStringList(newVal["values"].([]interface{})),
		}
	}

	return vals
}
