package scheduler

import (
	"context"
	"time"

	"dify-wp-sync/internal/dify"
	"dify-wp-sync/internal/logger"
	"dify-wp-sync/internal/sites"
	"dify-wp-sync/internal/wpcom"
)

func Start(ctx context.Context, interval time.Duration, sm *sites.Manager, difyCli *dify.DifyClient) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				syncAll(ctx, sm, difyCli)
			}
		}
	}()
}

func syncAll(ctx context.Context, sm *sites.Manager, difyCli *dify.DifyClient) {
	allSites, err := sm.ListSites(ctx)
	if err != nil {
		logger.Log.Errorf("Failed to list sites: %v", err)
		return
	}

	for _, sc := range allSites {
		err := wpcom.SyncSite(ctx, sc, difyCli)
		if err != nil {
			logger.Log.Errorf("Failed to sync site %s: %v", sc.SiteID, err)
			continue
		}
		if err := sm.UpdateSite(ctx, sc); err != nil {
			logger.Log.Errorf("Failed to update site %s after sync: %v", sc.SiteID, err)
		} else {
			logger.Log.Infof("Site %s synced successfully", sc.SiteID)
		}
	}
}
