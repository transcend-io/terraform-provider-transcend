package transcend

import "github.com/shurcooL/graphql"

func toStringList(raw interface{}) []graphql.String {
	if raw == nil {
		return []graphql.String{}
	}
	origs := raw.([]interface{})
	vals := make([]graphql.String, len(origs))
	for i, orig := range origs {
		vals[i] = graphql.String(orig.(string))
	}

	return vals
}

func toString(raw interface{}) graphql.String {
	if raw == nil {
		return ""
	}
	return graphql.String(raw.(string))
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
			Name:     graphql.String(newHead["name"].(string)),
			Value:    graphql.String(newHead["value"].(string)),
			IsSecret: graphql.Boolean(newHead["is_secret"].(bool)),
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

func flattenDataSiloBlockList(dataSilo DataSilo) []interface{} {
	owners := dataSilo.SubjectBlocklist
	ret := make([]interface{}, len(owners))
	for i, owner := range owners {
		ret[i] = owner.ID
	}
	return ret
}

func flattenOwners(dataSilo DataSilo) []interface{} {
	owners := dataSilo.Owners
	ret := make([]interface{}, len(owners))
	for i, owner := range owners {
		ret[i] = owner.Email
	}
	return ret
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
