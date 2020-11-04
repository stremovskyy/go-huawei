package internal

import (
	"net/url"
)

func SignURLWithApiKey(apiKey string) (string, error) {
	values := url.Values{}
	values.Add("key", apiKey)

	return values.Encode(), nil
}
