package wpcom

import (
	"dify-wp-sync/internal/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// WPClient interacts with the WordPress.com API.
type WPClient struct {
	AccessToken string
	SiteID      string
	httpClient  *http.Client
}

func NewWPClient(token, siteID string) *WPClient {
	return &WPClient{
		AccessToken: token,
		SiteID:      siteID,
		httpClient:  &http.Client{Timeout: 30 * time.Second},
	}
}

// GetPosts fetches the latest 20 posts ordered by modification time.
func (c *WPClient) GetPosts(modifiedAfter time.Time) ([]Post, error) {

	logger.Log.Infof("Fetching posts from site %s", c.SiteID)

	u := fmt.Sprintf("https://public-api.wordpress.com/rest/v1.1/sites/%s/posts", c.SiteID)
	params := url.Values{}
	params.Set("number", "20")
	params.Set("order_by", "modified")
	params.Set("order", "DESC")
	params.Set("fields", "posts(ID,date,modified,title,content),found")

	req, err := http.NewRequest("GET", u+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pr PostsResponse
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, err
	}

	var filtered []Post
	for _, p := range pr.Posts {
		if p.ModifiedTime().After(modifiedAfter) {
			filtered = append(filtered, p)
		}
	}

	return filtered, nil
}

// We will need to add pages (and other post types) here later.
