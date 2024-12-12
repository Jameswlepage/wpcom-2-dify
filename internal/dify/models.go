package dify

type CreateDatasetRequest struct {
	Name        string `json:"name"`
	Permission  string `json:"permission"`
	Description string `json:"description,omitempty"`
}

type DatasetResponse struct {
	ID string `json:"id"`
}

type CreateDocByTextRequest struct {
	Name              string      `json:"name"`
	Text              string      `json:"text"`
	IndexingTechnique string      `json:"indexing_technique"`
	ProcessRule       interface{} `json:"process_rule"`
}

type DocumentResponse struct {
	Document Document `json:"document"`
}

type Document struct {
	ID string `json:"id"`
}
