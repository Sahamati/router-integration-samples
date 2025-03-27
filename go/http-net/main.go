package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseURL             = "https://api.dev.sahamati.org.in/router"
	discoverURL         = baseURL + "/v2/Accounts/discover"
	recipientEntityID   = "FIP-SIMULATOR"
	defaultTxnID        = "f35761ac-4a18-11e8-96ff-0277a9fbfedc2"
	defaultJwsSignature = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	defaultSimulateRes  = "Ok"
	defaultAuthToken    = "token"
)

type SahamatiClient struct {
	httpClient *http.Client
	authToken  string
}

func NewSahamatiClient() *SahamatiClient {
	return &SahamatiClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		authToken:  defaultAuthToken,
	}
}

func (c *SahamatiClient) createDiscoverRequest() map[string]interface{} {
	return map[string]interface{}{
		"ver":       "2.0.0",
		"timestamp": time.Now().UTC().Format("2006-01-02T15:04:05.999Z07:00"),
		"txnid":     defaultTxnID,
		"Customer": map[string]interface{}{
			"id": "customer_identifier@AA_identifier",
			"Identifiers": []map[string]interface{}{
				{
					"category": "STRONG",
					"type":     "AADHAAR",
					"value":    "XXXXXXXXXXXXXXXX",
				},
			},
		},
		"FITypes": []string{"DEPOSIT"},
	}
}

func (c *SahamatiClient) setHeaders(req *http.Request) error {
	metadataHeader, err := encodeRequestMetadata(recipientEntityID)
	if err != nil {
		return fmt.Errorf("failed to encode metadata: %v", err)
	}

	req.Header.Set("x-jws-signature", defaultJwsSignature)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-simulate-res", defaultSimulateRes)
	req.Header.Set("Authorization", "Bearer "+c.authToken)
	req.Header.Set("x-request-meta", metadataHeader)

	return nil
}

func encodeRequestMetadata(entityID string) (string, error) {
	metadata := map[string]string{"recipient-id": entityID}
	jsonData, err := json.Marshal(metadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal metadata: %v", err)
	}
	return base64.StdEncoding.EncodeToString(jsonData), nil
}

func (c *SahamatiClient) DiscoverAccounts() ([]byte, error) {
	requestBody := c.createDiscoverRequest()

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, discoverURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	if err := c.setHeaders(req); err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return body, nil
}

func main() {
	httpClient := NewSahamatiClient()

	response, err := httpClient.DiscoverAccounts()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println(string(response))
}
