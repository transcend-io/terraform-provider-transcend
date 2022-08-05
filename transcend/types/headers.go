package types

import (
	graphql "github.com/hasura/go-graphql-client"
)

type Header struct {
	Name     graphql.String  `json:"name"`
	Value    graphql.String  `json:"value"`
	IsSecret graphql.Boolean `json:"isSecret"`
}

type CustomHeaderInput Header

func ToCustomHeaderInputList(origs []interface{}) []CustomHeaderInput {
	vals := make([]CustomHeaderInput, len(origs))
	for i, orig := range origs {
		newHead := orig.(map[string]interface{})

		vals[i] = CustomHeaderInput{
			Name:     graphql.String(newHead["name"].(string)),
			Value:    graphql.String(newHead["value"].(string)),
			IsSecret: graphql.Boolean(newHead["is_secret"].(bool)),
		}
	}

	return vals
}

func FlattenHeaders(headers *[]Header) []interface{} {
	ret := make([]interface{}, len(*headers))

	for i, header := range *headers {
		itemMap := make(map[string]interface{})
		itemMap["name"] = header.Name
		itemMap["value"] = header.Value
		itemMap["is_secret"] = header.IsSecret
		ret[i] = itemMap
	}

	return ret
}
