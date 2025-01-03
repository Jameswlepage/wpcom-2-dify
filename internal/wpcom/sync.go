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

	postTypes := siteCfg.PostTypes
	if len(postTypes) == 0 {
		postTypes = []string{"post"}
	}

	updatedSyncTime := siteCfg.LastSyncTime

	// Process each post type
	for _, postType := range postTypes {
		offset := 0
		limit := 100

		for {
			posts, hasMore, err := wp.GetPostsBatch(siteCfg.LastSyncTime, postType, offset, limit)
			if err != nil {
				return err
			}

			// Process this batch of posts
			for _, p := range posts {
				if p.Content == "" {
					logger.Log.Warnf("Post %d (%s) has empty content, skipping creation/update", p.ID, p.Title)
					continue
				}

				markdownContent := p.GetMarkdownContent()
				docID, exists := siteCfg.PostDocMapping[p.ID]

				if !exists {
					newDocID, err := difyClient.CreateDocumentByText(siteCfg.DifyDatasetID, p.Title, markdownContent)
					if err != nil {
						logger.Log.Errorf("Failed to create doc for post %d (%s): %v", p.ID, p.Title, err)
						continue
					}
					siteCfg.PostDocMapping[p.ID] = newDocID
					logger.Log.Infof("Created document %s for post %d (%s)", newDocID, p.ID, p.Title)
				} else {
					_, err := difyClient.UpdateDocumentByText(siteCfg.DifyDatasetID, docID, p.Title, markdownContent)
					if err != nil {
						logger.Log.Errorf("Failed to update doc %s for post %d (%s): %v", docID, p.ID, p.Title, err)
						continue
					}
					logger.Log.Infof("Updated document %s for post %d (%s)", docID, p.ID, p.Title)
				}

				if p.ModifiedTime().After(updatedSyncTime) {
					updatedSyncTime = p.ModifiedTime()
				}
			}

			if !hasMore {
				break
			}
			offset += limit
		}
	}

	siteCfg.LastSyncTime = updatedSyncTime
	return nil
}
