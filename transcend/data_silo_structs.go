package transcend

import (
	"github.com/shurcooL/graphql"
)

type DataSilo struct {
	ID    graphql.String
	Title graphql.String
	Link  graphql.String
}

type DataSiloRequest struct {
	Title  string `json:"title"`
	First  int    `json:"first"`
	Offset int    `json:"offset"`
}

type DataSiloResponse struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Link    string `json:"link"`
	Catalog struct {
		HasAvcFunctionality bool `json:"hasAvcFunctionality"`
	} `json:"catalog"`
}

// func getDataSilo(client genqlient.Client, title string, first, offset int) (*DataSiloResponse, error) {
// 	var retval DataSiloResponse
// 	var reqInput = DataSiloRequest{
// 		Title:  title,
// 		First:  first,
// 		Offset: offset,
// 	}
// 	err := client.MakeRequest(nil, "SchemaSyncDataSilos", `
// 		query SchemaSyncDataSilos(
// 			$title: String
// 			$first: Int!
// 			$offset: Int!,
// 		) {
// 			dataSilos(
// 			filterBy: { text: $title }
// 			first: $first
// 			offset: $offset
// 			) {
// 			nodes {
// 				id
// 				title
// 				link
// 				type
// 				catalog {
// 				hasAvcFunctionality
// 				}
// 			}
// 		}
// 	}
// 	`, &retval, &reqInput)

// 	return &retval, err
// }
