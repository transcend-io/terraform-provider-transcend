package transcend

import "github.com/shurcooL/graphql"

type DbIntegrationQuerySuggestionInput string

type RequestActionObjectResolver string

type DataPointSubDataPointInput struct {
	Name       graphql.String `json:"name"`
	Desciption graphql.String `json:"description"`
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
		}
	}

	return vals
}
