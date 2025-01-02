package wpcom

import (
	"context"
	"dify-wp-sync/internal/dify"
	"dify-wp-sync/internal/logger"
	"dify-wp-sync/internal/sites"
)

// SyncSite fetches posts of specified types updated since the site's last sync and
// either creates or updates corresponding documents in the Dify dataset.
func SyncSite(ctx context.Context, siteCfg *sites.SiteConfig, difyClient *dify.DifyClient) error {
	wp := NewWPClient(siteCfg.AccessToken, siteCfg.SiteID)

	// Determine which post types to synchronize
	postTypes := siteCfg.PostTypes
	if len(postTypes) == 0 {
		postTypes = []string{"post"} // Default to synchronizing posts only
	}

	posts, err := wp.GetPosts(siteCfg.LastSyncTime, postTypes)
	if err != nil {
		return err
	}
	if len(posts) == 0 {
		// No new or updated posts since last sync
		return nil
	}

	updatedSyncTime := siteCfg.LastSyncTime
	for _, p := range posts {
		// Skip if content is empty
		if p.Content == "" {
			logger.Log.Warnf("Post %d (%s) has empty content, skipping creation/update", p.ID, p.Title)
			continue
		}

		// Convert HTML to Markdown before sending to Dify
		markdownContent := p.GetMarkdownContent()

		docID, exists := siteCfg.PostDocMapping[p.ID]
		if !exists {
			// Create doc if it doesn't exist
			newDocID, err := difyClient.CreateDocumentByText(siteCfg.DifyDatasetID, p.Title, markdownContent)
			if err != nil {
				logger.Log.Errorf("Failed to create doc for post %d (%s): %v", p.ID, p.Title, err)
				continue
			}
			siteCfg.PostDocMapping[p.ID] = newDocID
			logger.Log.Infof("Created document %s for post %d (%s)", newDocID, p.ID, p.Title)
		} else {
			// Update doc if it exists
			_, err := difyClient.UpdateDocumentByText(siteCfg.DifyDatasetID, docID, p.Title, markdownContent)
			if err != nil {
				logger.Log.Errorf("Failed to update doc %s for post %d (%s): %v", docID, p.ID, p.Title, err)
				continue
			}
			logger.Log.Infof("Updated document %s for post %d (%s)", docID, p.ID, p.Title)
		}

		// Update sync time if post modified time is newer
		if p.ModifiedTime().After(updatedSyncTime) {
			updatedSyncTime = p.ModifiedTime()
		}
	}

	siteCfg.LastSyncTime = updatedSyncTime
	return nil
}
