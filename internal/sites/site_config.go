package sites

import (
	"time"
)

// SiteConfig represents the configuration for a WordPress site.
type SiteConfig struct {
	SiteID         string         `json:"site_id"`
	BlogURL        string         `json:"blog_url"`
	AccessToken    string         `json:"access_token"`
	DifyDatasetID  string         `json:"dify_dataset_id"`
	LastSyncTime   time.Time      `json:"last_sync_time"`
	PostDocMapping map[int]string `json:"post_doc_mapping"`
	PostTypes      []string       `json:"post_types"` // New field to specify post types to sync
}
