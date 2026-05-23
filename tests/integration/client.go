package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func makeRequest(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	client := http.Client{Timeout: 30 * time.Second}
	requestMessage, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	return client.Do(requestMessage)
}

func readResponseBody(response *http.Response, b any) error {
	defer response.Body.Close()
	if err := json.NewDecoder(response.Body).Decode(b); err != nil {
		return fmt.Errorf("decode response body: %w", err)
	}
	return nil
}
