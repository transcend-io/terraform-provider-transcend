package transcend

import (
	"net/http"

	"github.com/shurcooL/graphql"
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
	// token := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiToken})
	// client := oauth2.NewClient(context.Background(), token)
	client := &http.Client{Transport: &myTransport{apiToken: apiToken}}

	return &Client{
		graphql: graphql.NewClient(url, client),
		url:     url,
	}
}
