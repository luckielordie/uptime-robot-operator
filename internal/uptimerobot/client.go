package uptimerobot

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	apiKey string
}

func NewClient(apiKey string) Client {
	return Client{
		apiKey: apiKey,
	}
}

func makeApiRequest(ctx context.Context, client Client, methodName string, params map[string]string) ([]byte, error) {
	url := fmt.Sprintf("https://api.uptimerobot.com/v2/%s", methodName)
	payload := strings.NewReader(fmt.Sprintf("api_key=%s&format=json", client.apiKey))
	request, err := http.NewRequestWithContext(ctx, "POST", url, payload)
	if err != nil {
		return []byte{}, err
	}

	request.Header.Add("cache-control", "no-cache")
	request.Header.Add("content-type", "application/x-www-form-urlencoded")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return []byte{}, err
	}

	defer response.Body.Close()
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}

	return bodyBytes, nil
}
