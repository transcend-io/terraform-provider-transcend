package transcend

import (
	"net/http"

	graphql "github.com/hasura/go-graphql-client"
)

type backendTransport struct {
	apiToken string
}

func (t *backendTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer "+t.apiToken)
	return http.DefaultTransport.RoundTrip(req)
}

type sombraTransport struct {
	apiToken    string
	internalKey string
}

func (t *sombraTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer "+t.apiToken)
	if t.internalKey != "" {
		req.Header.Add("x-sombra-authorization", "Bearer "+t.internalKey)
	}
	return http.DefaultTransport.RoundTrip(req)
}


type Client struct {
   graphql         *graphql.Client
   sombraClient    *http.Client
   url             string
   internalSombraUrl string
}


func NewClient(url, apiToken string, internalKey string) *Client {
   return NewClientWithSombraUrl(url, apiToken, internalKey, "")
}

func NewClientWithSombraUrl(url, apiToken string, internalKey string, internalSombraUrl string) *Client {
   backendClient := &http.Client{Transport: &backendTransport{apiToken: apiToken}}
   sombraClient := &http.Client{Transport: &sombraTransport{apiToken: apiToken, internalKey: internalKey}}

   return &Client{
	   graphql:      graphql.NewClient(url, backendClient),
	   sombraClient: sombraClient,
	   url:          url,
	   internalSombraUrl: internalSombraUrl,
   }
}
