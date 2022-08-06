package transcend

import (
	"net/http"

	graphql "github.com/hasura/go-graphql-client"
)

type myTransport struct {
	apiToken string
}

func (t *myTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", t.apiToken)
	return http.DefaultTransport.RoundTrip(req)
}

type Client struct {
	graphql *graphql.Client
	url     string
}

func NewClient(url, apiToken string) *Client {
	apiToken = "Bearer " + apiToken
	client := &http.Client{Transport: &myTransport{apiToken: apiToken}}

	return &Client{
		graphql: graphql.NewClient(url, client),
		url:     url,
	}
}
