package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	oauthMgr := oauth.NewOAuthManager(cfg.ClientID, cfg.ClientSecret, cfg.RedirectURI)
	authHandler := &oauth.AuthHandler{Oauth: oauthMgr, SitesMgr: sitesMgr, DifyCli: difyClient}

	mux := http.NewServeMux()

	// Add the status endpoint:
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// This is your simple status page. You can add more info if needed.
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "System status: OK")
	})

	mux.HandleFunc("/oauth/callback", authHandler.HandleOAuthCallback)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		logger.Log.Infof("Starting server on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit
	logger.Log.Infof("Received signal %s, shutting down...", s)

	ctxShutdown, cancelShutdown := context.WithTimeout(ctx, 5*time.Second)
	defer cancelShutdown()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		logger.Log.Errorf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server exited")
}
