package transcend

import (
	"github.com/shurcooL/graphql"
)

type DataSilo struct {
	ID    graphql.String
	Title graphql.String
	Link  graphql.String
	Type  graphql.String
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
