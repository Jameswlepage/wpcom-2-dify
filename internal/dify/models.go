package dify

// CreateDatasetRequest is used for POST /datasets
type CreateDatasetRequest struct {
	Name        string `json:"name"`
	Permission  string `json:"permission"`
	Description string `json:"description,omitempty"`
}

// DatasetResponse is the response from creating a dataset.
type DatasetResponse struct {
	ID string `json:"id"`
}

// CreateDocByTextRequest is used for POST /datasets/:datasetID/document/create-by-text
type CreateDocByTextRequest struct {
	Name              string      `json:"name"`
	Text              string      `json:"text"`
	IndexingTechnique string      `json:"indexing_technique"`
	ProcessRule       interface{} `json:"process_rule"`
}

// DocumentResponse is the response when creating or updating a document.
type DocumentResponse struct {
	Document Document `json:"document"`
}

// Document encapsulates the ID of a single document.
type Document struct {
	ID string `json:"id"`
}

// ListDatasetsResponse is the JSON shape returned by GET /datasets.
type ListDatasetsResponse struct {
	Data    []DatasetInfo `json:"data"`
	HasMore bool          `json:"has_more"`
	Limit   int           `json:"limit"`
	Total   int           `json:"total"`
	Page    int           `json:"page"`
}

// DatasetInfo describes an individual dataset in the "data" array.
type DatasetInfo struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Permission     string `json:"permission"`
	DataSourceType string `json:"data_source_type"`
	// Add any other fields you might care about from the actual Dify JSON
	IndexingTechnique string `json:"indexing_technique"`
	AppCount          int    `json:"app_count"`
	DocumentCount     int    `json:"document_count"`
	WordCount         int    `json:"word_count"`
	CreatedBy         string `json:"created_by"`
	CreatedAt         int64  `json:"created_at"`
	UpdatedBy         string `json:"updated_by"`
	UpdatedAt         int64  `json:"updated_at"`
}
