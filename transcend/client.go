package transcend

import (
	"context"

	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

type Client struct {
	graphql *graphql.Client
	url     string
}

func NewClient(url, apiToken string) *Client {
	token := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiToken})
	client := oauth2.NewClient(context.Background(), token)

	return &Client{
		graphql: graphql.NewClient(url, client),
		url:     url,
	}
}
