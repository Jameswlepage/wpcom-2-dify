package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"dify-wp-sync/internal/config"
	"dify-wp-sync/internal/dify"
	"dify-wp-sync/internal/logger"
	"dify-wp-sync/internal/redisstore"
	"dify-wp-sync/internal/sites"
	"dify-wp-sync/internal/wpcom"
)

func main() {
	if len(os.Args) < 2 {
		printUsageAndExit()
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
	case "open-oauth":
		openOAuthPortal(cfg)
	case "force-sync-site":
		if len(os.Args) < 3 {
			fmt.Println("Usage: cli force-sync-site <site_id>")
			os.Exit(1)
		}
		siteID := os.Args[2]

		forceSyncSite(ctx, sitesMgr, siteID)
		syncSite(ctx, sitesMgr, difyClient, siteID)
	case "force-sync-doc":
		if len(os.Args) < 4 {
			fmt.Println("Usage: cli force-sync-doc <site_id> <post_id>")
			os.Exit(1)
		}

		siteID := os.Args[2]
		postIDStr := os.Args[3]
		postID, convErr := strconv.Atoi(postIDStr)
		if convErr != nil {
			fmt.Printf("Invalid post_id: %s\n", postIDStr)
			os.Exit(1)
		}
		forceSyncDoc(ctx, sitesMgr, siteID, postID)
	case "set-post-types":
		if len(os.Args) < 4 {
			fmt.Println("Usage: cli set-post-types <site_id> <post_types_comma_separated>")
			os.Exit(1)
		}
		siteID := os.Args[2]
		postTypesStr := os.Args[3]
		setSitePostTypes(ctx, sitesMgr, siteID, postTypesStr)
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		os.Exit(1)
	}
}

func printUsageAndExit() {
	fmt.Println("Usage: cli <command> [args...]")
	fmt.Println("Commands:")
	fmt.Println("  list-sites")
	fmt.Println("  sync-site <site_id>")
	fmt.Println("  sync-all-sites")
	fmt.Println("  open-oauth")
	fmt.Println("  force-sync-site <site_id>")
	fmt.Println("  force-sync-doc <site_id> <post_id>")
	fmt.Println("  set-post-types <site_id> <post_types_comma_separated>")
	os.Exit(1)
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
		fmt.Printf("- SiteID: %s, BlogURL: %s, LastSync: %s, PostTypes: %v\n", s.SiteID, s.BlogURL, s.LastSyncTime, s.PostTypes)
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

func openOAuthPortal(cfg *config.Config) {
	oauthURL := fmt.Sprintf("https://public-api.wordpress.com/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code",
		cfg.ClientID, url.QueryEscape(cfg.RedirectURI))

	fmt.Println("Open the following URL in your browser to authorize your site:")
	fmt.Println(oauthURL)
}

func forceSyncSite(ctx context.Context, sm *sites.Manager, siteID string) {
	sc, err := sm.GetSite(ctx, siteID)
	if err != nil {
		logger.Log.Errorf("Failed to get site %s for force-sync: %v", siteID, err)
		os.Exit(1)
	}
	sc.PostDocMapping = make(map[int]string)
	sc.LastSyncTime = time.Time{}
	if err := sm.UpdateSite(ctx, sc); err != nil {
		logger.Log.Errorf("Failed to update site %s for force-sync: %v", siteID, err)
		os.Exit(1)
	}
	fmt.Printf("Site %s has been reset. The next sync will recreate all documents.\n", siteID)
}

func forceSyncDoc(ctx context.Context, sm *sites.Manager, siteID string, postID int) {
	sc, err := sm.GetSite(ctx, siteID)
	if err != nil {
		logger.Log.Errorf("Failed to get site %s for force-sync-doc: %v", siteID, err)
		os.Exit(1)
	}
	delete(sc.PostDocMapping, postID)
	if err := sm.UpdateSite(ctx, sc); err != nil {
		logger.Log.Errorf("Failed to update site %s after removing doc mapping for post %d: %v", siteID, postID, err)
		os.Exit(1)
	}
	fmt.Printf("Document mapping for post %d on site %s removed. Run 'sync-site %s' again to recreate.\n", postID, siteID, siteID)
}

func setSitePostTypes(ctx context.Context, sm *sites.Manager, siteID, postTypesStr string) {
	sc, err := sm.GetSite(ctx, siteID)
	if err != nil {
		logger.Log.Errorf("Failed to get site %s for setting post types: %v", siteID, err)
		os.Exit(1)
	}
	postTypes := strings.Split(postTypesStr, ",")
	sc.PostTypes = postTypes
	if err := sm.UpdateSite(ctx, sc); err != nil {
		logger.Log.Errorf("Failed to update site %s after setting post types: %v", siteID, err)
		os.Exit(1)
	}
	fmt.Printf("Post types for site %s updated to: %v\n", siteID, postTypes)
}
