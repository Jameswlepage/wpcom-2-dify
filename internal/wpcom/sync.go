package wpcom

import (
	"context"

	"dify-wp-sync/internal/dify"
	"dify-wp-sync/internal/logger"
	"dify-wp-sync/internal/sites"
)

func SyncSite(ctx context.Context, siteCfg *sites.SiteConfig, difyClient *dify.DifyClient) error {
	wp := NewWPClient(siteCfg.AccessToken, siteCfg.SiteID)
	posts, err := wp.GetPosts(siteCfg.LastSyncTime)
	if err != nil {
		return err
	}
	if len(posts) == 0 {
		return nil
	}

	updatedSyncTime := siteCfg.LastSyncTime
	for _, p := range posts {
		docID, exists := siteCfg.PostDocMapping[p.ID]
		if !exists {
			// create doc
			newDocID, err := difyClient.CreateDocumentByText(siteCfg.DifyDatasetID, p.Title, p.ContentPlain)
			if err != nil {
				logger.Log.Errorf("Failed to create doc for post %d: %v", p.ID, err)
				continue
			}
			siteCfg.PostDocMapping[p.ID] = newDocID
			logger.Log.Infof("Created document %s for post %d (%s)", newDocID, p.ID, p.Title)
		} else {
			// update doc
			_, err := difyClient.UpdateDocumentByText(siteCfg.DifyDatasetID, docID, p.Title, p.ContentPlain)
			if err != nil {
				logger.Log.Errorf("Failed to update doc %s for post %d: %v", docID, p.ID, err)
				continue
			}
			logger.Log.Infof("Updated document %s for post %d (%s)", docID, p.ID, p.Title)
		}

		if p.ModifiedTime().After(updatedSyncTime) {
			updatedSyncTime = p.ModifiedTime()
		}
	}
	siteCfg.LastSyncTime = updatedSyncTime
	return nil
}
