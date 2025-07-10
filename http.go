package openxbl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *Client) makeRequest(method string, endpoint string, body any, object any) ([]byte, error) {
	var err error
	var requestBytes, responseBytes []byte

	if body != nil {
		requestBytes, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshalling body: %w", err)
		}
	}

	req, err := http.NewRequest(strings.ToUpper(method), url+endpoint, bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Authorization", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.Body != nil {
		defer resp.Body.Close()

		responseBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %w", err)
		}
	}

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	// Unmarshal to object if one is provided.
	if object != nil && len(responseBytes) > 0 {
		if err = json.Unmarshal(responseBytes, &object); err != nil {
			return nil, fmt.Errorf("error unmarshaling response: %w", err)
		}
	}

	return responseBytes, nil
}
