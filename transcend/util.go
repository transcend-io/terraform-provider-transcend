package transcend

import (
	"github.com/shurcooL/graphql"
	"github.com/transcend-io/terraform-provider-transcend/transcend/types"
)

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

func toCustomHeaderInputList(origs []interface{}) []types.CustomHeaderInput {
	vals := make([]types.CustomHeaderInput, len(origs))
	for i, orig := range origs {
		newHead := orig.(map[string]interface{})

		vals[i] = types.CustomHeaderInput{
			Name:     graphql.String(newHead["name"].(string)),
			Value:    graphql.String(newHead["value"].(string)),
			IsSecret: graphql.Boolean(newHead["is_secret"].(bool)),
		}
	}

	return vals
}

func toRequestActionList(origs []interface{}) []types.RequestAction {
	vals := make([]types.RequestAction, len(origs))
	for i, orig := range origs {
		vals[i] = types.RequestAction(orig.(string))
	}

	return vals
}

func flattenDataSiloBlockList(dataSilo types.DataSilo) []interface{} {
	owners := dataSilo.SubjectBlocklist
	ret := make([]interface{}, len(owners))
	for i, owner := range owners {
		ret[i] = owner.ID
	}
	return ret
}

func flattenOwners(dataSilo types.DataSilo) []interface{} {
	owners := dataSilo.Owners
	ret := make([]interface{}, len(owners))
	for i, owner := range owners {
		ret[i] = owner.Email
	}
	return ret
}

func flattenHeaders(headers *[]types.Header) []interface{} {
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

func toRequestActionObjectResolverList(origs []interface{}) []types.RequestActionObjectResolver {
	vals := make([]types.RequestActionObjectResolver, len(origs))
	for i, orig := range origs {
		vals[i] = types.RequestActionObjectResolver(orig.(string))
	}

	return vals
}

func toDataPointSubDataPointInputList(origs []interface{}) []types.DataPointSubDataPointInput {
	vals := make([]types.DataPointSubDataPointInput, len(origs))
	for i, orig := range origs {
		newVal := orig.(map[string]interface{})
		vals[i] = types.DataPointSubDataPointInput{
			graphql.String(newVal["name"].(string)),
			graphql.String(newVal["description"].(string)),
			toDataSubCategoryInputList(newVal["categories"].([]interface{})),
			toPurposeSubCategoryInputList(newVal["purposes"].([]interface{})),
			toAttributeInputList(newVal["attributes"].([]interface{})),
		}
	}

	return vals
}

func toDataSubCategoryInputList(origs []interface{}) []types.DataSubCategoryInput {
	vals := make([]types.DataSubCategoryInput, len(origs))
	for i, orig := range origs {
		newVal := orig.(map[string]interface{})
		vals[i] = types.DataSubCategoryInput{
			graphql.String(newVal["name"].(string)),
			types.DataCategoryType(newVal["category"].(string)),
		}
	}

	return vals
}

func toPurposeSubCategoryInputList(origs []interface{}) []types.PurposeSubCategoryInput {
	vals := make([]types.PurposeSubCategoryInput, len(origs))
	for i, orig := range origs {
		newVal := orig.(map[string]interface{})
		vals[i] = types.PurposeSubCategoryInput{
			graphql.String(newVal["name"].(string)),
			types.ProcessingPurpose(newVal["purpose"].(string)),
		}
	}

	return vals
}

func toAttributeInputList(origs []interface{}) []types.AttributeInput {
	vals := make([]types.AttributeInput, len(origs))
	for i, orig := range origs {
		newVal := orig.(map[string]interface{})
		vals[i] = types.AttributeInput{
			graphql.String(newVal["key"].(string)),
			toStringList(newVal["values"].([]interface{})),
		}
	}

	return vals
}
