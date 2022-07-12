package transcend

import (
	"github.com/shurcooL/graphql"
)

type APIKey struct {
	ID     graphql.String
	Title  graphql.String
	Scopes []struct {
		Type graphql.String
	}
	DataSilos []DataSilo
}
