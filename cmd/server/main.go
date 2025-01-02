package main

import (
	"fmt"
	"net/http"
	"os"

	"dify-wp-sync/internal/config"
	"dify-wp-sync/internal/dify"
	"dify-wp-sync/internal/logger"
	"dify-wp-sync/internal/oauth"
	"dify-wp-sync/internal/redisstore"
	"dify-wp-sync/internal/sites"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Log.Fatalf("Error loading config: %v", err)
	}

	store := redisstore.New(cfg.RedisAddr, cfg.RedisPwd, cfg.RedisDB)
	sitesMgr := sites.NewManager(store)
	difyClient := dify.NewDifyClient(cfg.DifyToken, cfg.DifyBaseURL)
	oauthManager := oauth.NewOAuthManager(cfg.ClientID, cfg.ClientSecret, cfg.RedirectURI)
	authHandler := &oauth.AuthHandler{
		Oauth:    oauthManager,
		SitesMgr: sitesMgr,
		DifyCli:  difyClient,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "System status: OK")
	})
	http.HandleFunc("/oauth/callback", authHandler.HandleOAuthCallback)

	logger.Log.Infof("Starting server on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		logger.Log.Fatalf("Server failed to start: %v", err)
		os.Exit(1)
	}
}
