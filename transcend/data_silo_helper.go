package transcend

import (
	"github.com/shurcooL/graphql"
)

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

func flattenItems(items *[]DataSilo) []interface{} {
	if items == nil {
		return make([]interface{}, 0)
	}

	ret := make([]interface{}, len(*items))

	for i, item := range *items {
		itemMap := make(map[string]interface{})
		itemMap["id"] = item.ID
		itemMap["title"] = item.Title
		itemMap["link"] = item.Link
		itemMap["type"] = item.Type

		ret[i] = itemMap
	}

	return ret
}
