package transcend

import (
	"github.com/shurcooL/graphql"
)

type DataSilo struct {
	ID    graphql.String
	Title graphql.String
	Link  graphql.String
}

type InputDataSilo struct {
	Name               graphql.String
	Title              graphql.String
	Description        graphql.String
	URL                graphql.String
	NotifyEmailAddress graphql.String
	IsLive             graphql.Boolean
	APIKeyID           graphql.String
}
