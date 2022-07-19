package transcend

import "github.com/shurcooL/graphql"

type DataPoint struct {
	Title struct {
		DefaultMessage graphql.String
	}
	Description struct {
		DefaultMessage graphql.String
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

func toDbIntegrationQuerySuggestionInputList(origs []interface{}) []DbIntegrationQuerySuggestionInput {
	vals := make([]DbIntegrationQuerySuggestionInput, len(origs))
	for i, orig := range origs {
		vals[i] = DbIntegrationQuerySuggestionInput(orig.(string))
	}

	return vals
}

func toRequestActionObjectResolverList(origs []interface{}) []RequestActionObjectResolver {
	vals := make([]RequestActionObjectResolver, len(origs))
	for i, orig := range origs {
		vals[i] = RequestActionObjectResolver(orig.(string))
	}

	return vals
}

func toDataPointSubDataPointInputList(origs []interface{}) []DataPointSubDataPointInput {
	vals := make([]DataPointSubDataPointInput, len(origs))
	for i, orig := range origs {
		newVal := orig.(map[string]interface{})
		vals[i] = DataPointSubDataPointInput{
			graphql.String(newVal["name"].(string)),
			graphql.String(newVal["description"].(string)),
			toDataSubCategoryInputList(newVal["categories"].([]interface{})),
			toPurposeSubCategoryInputList(newVal["purposes"].([]interface{})),
			toAttributeInputList(newVal["attributes"].([]interface{})),
		}
	}

	return vals
}

func toDataSubCategoryInputList(origs []interface{}) []DataSubCategoryInput {
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

func toPurposeSubCategoryInputList(origs []interface{}) []PurposeSubCategoryInput {
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

func toAttributeInputList(origs []interface{}) []AttributeInput {
	vals := make([]AttributeInput, len(origs))
	for i, orig := range origs {
		newVal := orig.(map[string]interface{})
		vals[i] = AttributeInput{
			graphql.String(newVal["key"].(string)),
			toStringList(newVal["values"].([]interface{})),
		}
	}

	return vals
}
