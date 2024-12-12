package sites

import "time"

// SiteConfig holds the configuration for a single WordPress site.
type SiteConfig struct {
	SiteID         string         `json:"site_id"` // WordPress.com blog_id or Jetpack site_id
	AccessToken    string         `json:"access_token"`
	BlogURL        string         `json:"blog_url"`
	DifyDatasetID  string         `json:"dify_dataset_id"`
	LastSyncTime   time.Time      `json:"last_sync_time"`
	PostDocMapping map[int]string `json:"post_doc_mapping"` // postID -> dify docID
}
