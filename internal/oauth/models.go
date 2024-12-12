package oauth

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	BlogID      string `json:"blog_id"`
	BlogURL     string `json:"blog_url"`
	TokenType   string `json:"token_type"`
}
