package transcend

import "github.com/shurcooL/graphql"

func toStringList(origs []interface{}) []graphql.String {
	vals := make([]graphql.String, len(origs))
	for i, orig := range origs {
		vals[i] = graphql.String(orig.(string))
	}

	return vals
}

func toIDList(origs []interface{}) []graphql.ID {
	vals := make([]graphql.ID, len(origs))
	for i, orig := range origs {
		vals[i] = graphql.ID(orig.(string))
	}

	return vals
}

func toCustomHeaderInputList(origs []interface{}) []CustomHeaderInput {
	vals := make([]CustomHeaderInput, len(origs))
	for i, orig := range origs {
		newHead := orig.(map[string]interface{})

		vals[i] = CustomHeaderInput{
			graphql.String(newHead["name"].(string)),
			graphql.String(newHead["value"].(string)),
			graphql.Boolean(newHead["is_secret"].(bool)),
		}
	}

	return vals
}

func toRequestActionList(origs []interface{}) []RequestAction {
	vals := make([]RequestAction, len(origs))
	for i, orig := range origs {
		vals[i] = RequestAction(orig.(string))
	}

	return vals
}

func flattenHeaders(headers *[]Header) []interface{} {
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
