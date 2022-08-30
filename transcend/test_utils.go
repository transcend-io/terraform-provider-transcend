package transcend

import "os"

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getTestClient() *Client {
	return NewClient(getEnv("TRANSCEND_URL", "https://yo.com:4001/graphql"), os.Getenv("TRANSCEND_KEY"), "")
}
