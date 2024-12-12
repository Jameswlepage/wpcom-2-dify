package oauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type OAuthManager struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

func NewOAuthManager(clientID, clientSecret, redirectURI string) *OAuthManager {
	return &OAuthManager{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
	}
}

func (o *OAuthManager) ExchangeCodeForToken(code string) (*TokenResponse, error) {
	// Exchange the code for an access token
	form := url.Values{}
	form.Set("client_id", o.ClientID)
	form.Set("client_secret", o.ClientSecret)
	form.Set("redirect_uri", o.RedirectURI)
	form.Set("code", code)
	form.Set("grant_type", "authorization_code")

	req, err := http.NewRequest("POST", "https://public-api.wordpress.com/oauth2/token", bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code %d from token endpoint", resp.StatusCode)
	}

	var tr TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return nil, err
	}
	return &tr, nil
}
