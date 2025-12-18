package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultBaseURL = "https://api.porkbun.com/api/json/v3"
)

// Client is the Porkbun API client
type Client struct {
	baseURL      string
	apiKey       string
	secretAPIKey string
	httpClient   *http.Client
}

// NewClient creates a new Porkbun API client
func NewClient(apiKey, secretAPIKey string) *Client {
	return &Client{
		baseURL:      defaultBaseURL,
		apiKey:       apiKey,
		secretAPIKey: secretAPIKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// authRequest is the base request with authentication
type authRequest struct {
	SecretAPIKey string `json:"secretapikey"`
	APIKey       string `json:"apikey"`
}

// APIResponse is the base response from the API
type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// DNSRecord represents a DNS record
type DNSRecord struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
	TTL     string `json:"ttl"`
	Prio    string `json:"prio"`
	Notes   string `json:"notes"`
}

// CreateDNSRecordRequest is the request to create a DNS record
type CreateDNSRecordRequest struct {
	authRequest
	Name    string `json:"name,omitempty"`
	Type    string `json:"type"`
	Content string `json:"content"`
	TTL     string `json:"ttl,omitempty"`
	Prio    string `json:"prio,omitempty"`
	Notes   string `json:"notes,omitempty"`
}

// CreateDNSRecordResponse is the response from creating a DNS record
type CreateDNSRecordResponse struct {
	APIResponse
	ID int64 `json:"id"`
}

// EditDNSRecordRequest is the request to edit a DNS record
type EditDNSRecordRequest struct {
	authRequest
	Name    string `json:"name,omitempty"`
	Type    string `json:"type"`
	Content string `json:"content"`
	TTL     string `json:"ttl,omitempty"`
	Prio    string `json:"prio,omitempty"`
	Notes   string `json:"notes,omitempty"`
}

// RetrieveDNSRecordsResponse is the response from retrieving DNS records
type RetrieveDNSRecordsResponse struct {
	APIResponse
	Records []DNSRecord `json:"records"`
}

// doRequest performs an HTTP request to the Porkbun API
func (c *Client) doRequest(method, endpoint string, body interface{}) ([]byte, error) {
	url := c.baseURL + endpoint

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// Ping tests the API connection
func (c *Client) Ping() error {
	req := authRequest{
		SecretAPIKey: c.secretAPIKey,
		APIKey:       c.apiKey,
	}

	respBody, err := c.doRequest("POST", "/ping", req)
	if err != nil {
		return err
	}

	var resp APIResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if resp.Status != "SUCCESS" {
		return fmt.Errorf("API ping failed: %s", resp.Message)
	}

	return nil
}

// CreateDNSRecord creates a new DNS record
func (c *Client) CreateDNSRecord(domain string, record CreateDNSRecordRequest) (string, error) {
	record.SecretAPIKey = c.secretAPIKey
	record.APIKey = c.apiKey

	respBody, err := c.doRequest("POST", "/dns/create/"+domain, record)
	if err != nil {
		return "", err
	}

	var resp CreateDNSRecordResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if resp.Status != "SUCCESS" {
		return "", fmt.Errorf("failed to create DNS record: %s", resp.Message)
	}

	return fmt.Sprintf("%d", resp.ID), nil
}

// GetDNSRecord retrieves a specific DNS record by ID
func (c *Client) GetDNSRecord(domain, recordID string) (*DNSRecord, error) {
	req := authRequest{
		SecretAPIKey: c.secretAPIKey,
		APIKey:       c.apiKey,
	}

	respBody, err := c.doRequest("POST", "/dns/retrieve/"+domain+"/"+recordID, req)
	if err != nil {
		return nil, err
	}

	var resp RetrieveDNSRecordsResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if resp.Status != "SUCCESS" {
		return nil, fmt.Errorf("failed to retrieve DNS record: %s", resp.Message)
	}

	if len(resp.Records) == 0 {
		return nil, fmt.Errorf("DNS record not found")
	}

	return &resp.Records[0], nil
}

// EditDNSRecord updates an existing DNS record
func (c *Client) EditDNSRecord(domain, recordID string, record EditDNSRecordRequest) error {
	record.SecretAPIKey = c.secretAPIKey
	record.APIKey = c.apiKey

	respBody, err := c.doRequest("POST", "/dns/edit/"+domain+"/"+recordID, record)
	if err != nil {
		return err
	}

	var resp APIResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if resp.Status != "SUCCESS" {
		return fmt.Errorf("failed to edit DNS record: %s", resp.Message)
	}

	return nil
}

// DeleteDNSRecord deletes a DNS record
func (c *Client) DeleteDNSRecord(domain, recordID string) error {
	req := authRequest{
		SecretAPIKey: c.secretAPIKey,
		APIKey:       c.apiKey,
	}

	respBody, err := c.doRequest("POST", "/dns/delete/"+domain+"/"+recordID, req)
	if err != nil {
		return err
	}

	var resp APIResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if resp.Status != "SUCCESS" {
		return fmt.Errorf("failed to delete DNS record: %s", resp.Message)
	}

	return nil
}
