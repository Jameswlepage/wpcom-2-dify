package dify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DifyClient struct {
	token   string
	baseURL string
	client  *http.Client
}

func NewDifyClient(token, baseURL string) *DifyClient {
	return &DifyClient{
		token:   token,
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

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

func (d *DifyClient) CreateDocumentByText(datasetID, name, text string) (string, error) {
	reqBody := CreateDocByTextRequest{
		Name:              name,
		Text:              text,
		IndexingTechnique: "high_quality",
		ProcessRule:       map[string]string{"mode": "automatic"},
	}
	b, _ := json.Marshal(reqBody)
	url := fmt.Sprintf("%s/datasets/%s/document/create-by-text", d.baseURL, datasetID)
	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
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

func (d *DifyClient) UpdateDocumentByText(datasetID, docID, name, text string) (string, error) {
	reqBody := map[string]string{
		"name": name,
		"text": text,
	}
	b, _ := json.Marshal(reqBody)
	url := fmt.Sprintf("%s/datasets/%s/documents/%s/update_by_text", d.baseURL, datasetID, docID)
	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
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
