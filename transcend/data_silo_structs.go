package transcend

import (
	"github.com/shurcooL/graphql"
)

type DataSilo struct {
	ID    graphql.String
	Title graphql.String
	Link  graphql.String
}
