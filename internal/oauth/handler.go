package oauth

import (
	"fmt"
	"net/http"

	"dify-wp-sync/internal/dify"
	"dify-wp-sync/internal/logger"
	"dify-wp-sync/internal/sites"
)

// AuthHandler handles the OAuth callback from WordPress.com.
// When a user authorizes your application, WordPress.com redirects here with a code.
// AuthHandler exchanges the code for a token, creates a Dify dataset, and stores the site config.
type AuthHandler struct {
	Oauth    *OAuthManager
	SitesMgr *sites.Manager
	DifyCli  *dify.DifyClient
}

// HandleOAuthCallback processes the authorization code returned by WordPress.com,
// exchanges it for an access token, and registers the site in Redis with a corresponding Dify dataset.
func (ah *AuthHandler) HandleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if errVal := r.URL.Query().Get("error"); errVal != "" {
		http.Error(w, "User denied access", http.StatusUnauthorized)
		return
	}
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No code provided", http.StatusBadRequest)
		return
	}

	tr, err := ah.Oauth.ExchangeCodeForToken(code)
	if err != nil {
		logger.Log.Errorf("Error exchanging code for token: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	datasetID, err := ah.DifyCli.CreateDataset(tr.BlogURL)
	if err != nil {
		logger.Log.Errorf("Failed to create Dify dataset: %v", err)
		http.Error(w, "Failed to create dataset", http.StatusInternalServerError)
		return
	}

	sc := &sites.SiteConfig{
		SiteID:        tr.BlogID,
		AccessToken:   tr.AccessToken,
		BlogURL:       tr.BlogURL,
		DifyDatasetID: datasetID,
	}

	if err := ah.SitesMgr.AddSite(ctx, sc); err != nil {
		logger.Log.Errorf("Failed to store site config: %v", err)
		http.Error(w, "Failed to store site config", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Site connected: %s (dataset: %s)", tr.BlogURL, datasetID)
}
