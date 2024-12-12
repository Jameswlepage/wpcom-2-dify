package main

import (
	"context"
	"fmt"
	"os"

	"dify-wp-sync/internal/config"
	"dify-wp-sync/internal/dify"
	"dify-wp-sync/internal/logger"
	"dify-wp-sync/internal/redisstore"
	"dify-wp-sync/internal/sites"
	"dify-wp-sync/internal/wpcom"
)

// This CLI tool provides commands to manage WordPress sites and synchronize them to Dify.
// Commands:
//
//	list-sites         - Lists all registered sites
//	sync-site <id>     - Syncs a single site by its SiteID
//	sync-all-sites     - Syncs all registered sites
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cli <command> [args...]")
		fmt.Println("Commands:")
		fmt.Println("  list-sites")
		fmt.Println("  sync-site <site_id>")
		fmt.Println("  sync-all-sites")
		os.Exit(1)
	}

	cmd := os.Args[1]

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Log.Fatalf("Error loading config: %v", err)
	}
	store := redisstore.New(cfg.RedisAddr, cfg.RedisPwd, cfg.RedisDB)
	sitesMgr := sites.NewManager(store)
	difyClient := dify.NewDifyClient(cfg.DifyToken, cfg.DifyBaseURL)
	ctx := context.Background()

	switch cmd {
	case "list-sites":
		listSites(ctx, sitesMgr)
	case "sync-site":
		if len(os.Args) < 3 {
			fmt.Println("Usage: cli sync-site <site_id>")
			os.Exit(1)
		}
		siteID := os.Args[2]
		syncSite(ctx, sitesMgr, difyClient, siteID)
	case "sync-all-sites":
		syncAllSites(ctx, sitesMgr, difyClient)
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		os.Exit(1)
	}
}

func listSites(ctx context.Context, sm *sites.Manager) {
	allSites, err := sm.ListSites(ctx)
	if err != nil {
		logger.Log.Errorf("Failed to list sites: %v", err)
		os.Exit(1)
	}
	if len(allSites) == 0 {
		fmt.Println("No sites found.")
		return
	}
	fmt.Println("Registered Sites:")
	for _, s := range allSites {
		fmt.Printf("- SiteID: %s, BlogURL: %s, LastSync: %s\n", s.SiteID, s.BlogURL, s.LastSyncTime)
	}
}

func syncSite(ctx context.Context, sm *sites.Manager, difyCli *dify.DifyClient, siteID string) {
	sc, err := sm.GetSite(ctx, siteID)
	if err != nil {
		logger.Log.Errorf("Failed to get site %s: %v", siteID, err)
		os.Exit(1)
	}
	err = wpcom.SyncSite(ctx, sc, difyCli)
	if err != nil {
		logger.Log.Errorf("Failed to sync site %s: %v", siteID, err)
		os.Exit(1)
	}
	if err := sm.UpdateSite(ctx, sc); err != nil {
		logger.Log.Errorf("Failed to update site %s after sync: %v", siteID, err)
		os.Exit(1)
	}
	fmt.Printf("Site %s synced successfully.\n", siteID)
}

func syncAllSites(ctx context.Context, sm *sites.Manager, difyCli *dify.DifyClient) {
	allSites, err := sm.ListSites(ctx)
	if err != nil {
		logger.Log.Errorf("Failed to list sites: %v", err)
		os.Exit(1)
	}
	if len(allSites) == 0 {
		fmt.Println("No sites to sync.")
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
			continue
		}
		fmt.Printf("Site %s synced successfully.\n", sc.SiteID)
	}
}
