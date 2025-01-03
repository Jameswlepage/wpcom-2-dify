package dify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// DifyClient provides methods for interacting with the Dify API.
type DifyClient struct {
	token   string
	baseURL string
	client  *http.Client
}

func NewDifyClient(token, baseURL string) *DifyClient {
	return &DifyClient{
		token:   token,
		baseURL: baseURL,
		client:  &http.Client{Timeout: 15 * time.Second},
	}
}

// ListDatasets calls GET /datasets?page=PAGE&limit=LIMIT
// and returns the parsed JSON response.
func (d *DifyClient) ListDatasets(page, limit int) (*ListDatasetsResponse, error) {
	url := fmt.Sprintf("%s/datasets?page=%d&limit=%d", d.baseURL, page, limit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+d.token)

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("list datasets returned status code %d", resp.StatusCode)
	}

	var listResp ListDatasetsResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("failed to decode ListDatasets response: %w", err)
	}
	return &listResp, nil
}

// DatasetExists enumerates all pages of datasets, checking if datasetID is present.
func (d *DifyClient) DatasetExists(datasetID string) (bool, error) {
	const pageSize = 20
	page := 1

	for {
		listResp, err := d.ListDatasets(page, pageSize)
		if err != nil {
			return false, fmt.Errorf("failed to list datasets (page=%d): %w", page, err)
		}
		// Check each dataset returned for a matching ID
		for _, ds := range listResp.Data {
			if ds.ID == datasetID {
				return true, nil
			}
		}
		// If we have exhausted all data, break
		if !listResp.HasMore || (page*pageSize >= listResp.Total) {
			break
		}
		page++
	}
	return false, nil
}

// CreateDataset calls Dify to create a new (empty) dataset with "only_me" permission.
func (d *DifyClient) CreateDataset(name string) (string, error) {
	reqBody := CreateDatasetRequest{
		Name:       name,
		Permission: "only_me",
	}
	b, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", d.baseURL+"/datasets", bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+d.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("failed to create dataset, status %d", resp.StatusCode)
	}
	var dsr DatasetResponse
	if err := json.NewDecoder(resp.Body).Decode(&dsr); err != nil {
		return "", err
	}
	return dsr.ID, nil
}

// CreateDocumentByText creates a new document in the specified dataset.
func (d *DifyClient) CreateDocumentByText(datasetID, name, text string) (string, error) {
	reqBody := CreateDocByTextRequest{
		Name:              name,
		Text:              text,
		IndexingTechnique: "high_quality",
		ProcessRule:       map[string]string{"mode": "automatic"},
	}
	b, _ := json.Marshal(reqBody)
	fullURL := fmt.Sprintf("%s/datasets/%s/document/create-by-text", d.baseURL, datasetID)
	req, err := http.NewRequest("POST", fullURL, bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+d.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("failed to create document, status %d", resp.StatusCode)
	}
	var dr DocumentResponse
	if err := json.NewDecoder(resp.Body).Decode(&dr); err != nil {
		return "", err
	}
	return dr.Document.ID, nil
}

// UpdateDocumentByText updates an existing Dify document with new text content.
func (d *DifyClient) UpdateDocumentByText(datasetID, docID, name, text string) (string, error) {
	reqBody := map[string]string{
		"name": name,
		"text": text,
	}
	b, _ := json.Marshal(reqBody)
	fullURL := fmt.Sprintf("%s/datasets/%s/documents/%s/update_by_text", d.baseURL, datasetID, docID)
	req, err := http.NewRequest("POST", fullURL, bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+d.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("failed to update document, status %d", resp.StatusCode)
	}
	var dr DocumentResponse
	if err := json.NewDecoder(resp.Body).Decode(&dr); err != nil {
		return "", err
	}
	return dr.Document.ID, nil
}
