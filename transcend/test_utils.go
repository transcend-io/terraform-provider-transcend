package transcend

import (
	"net/url"
	"os"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getTestClient() *Client {
	backendUrl := getEnv("TRANSCEND_URL", "https://yo.com:4001/")
	graphQlUrl, _ := url.JoinPath(backendUrl, "/graphql")

	return NewClient(graphQlUrl, os.Getenv("TRANSCEND_KEY"), "")
}
