package wpcom

import (
	"dify-wp-sync/internal/logger"
	"encoding/json"
	"fmt"
	"io"
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

// GetPosts fetches posts of specified types updated since modifiedAfter.
func (c *WPClient) GetPosts(modifiedAfter time.Time, postTypes []string) ([]Post, error) {
	var allPosts []Post
	for _, postType := range postTypes {
		logger.Log.Infof("Fetching posts of type '%s' from site %s", postType, c.SiteID)
		u := fmt.Sprintf("https://public-api.wordpress.com/rest/v1.1/sites/%s/posts", c.SiteID)
		params := url.Values{}
		params.Set("number", "100")
		params.Set("order_by", "modified")
		params.Set("order", "DESC")
		params.Set("fields", "ID,date,modified,title,content,type")
		params.Set("type", postType)

		req, err := http.NewRequest("GET", u+"?"+params.Encode(), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+c.AccessToken)
		req.Header.Set("User-Agent", "Dify-WP-Sync/1.0")
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		logger.Log.Infof("Received response status: %d from WordPress API", resp.StatusCode)

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %v", err)
		}
		logger.Log.Debugf("Raw response body: %s", string(bodyBytes))

		if resp.StatusCode != 200 {
			logger.Log.Errorf("Non-200 response: %d, body: %s", resp.StatusCode, string(bodyBytes))
			return nil, fmt.Errorf("unexpected status code %d from WordPress API", resp.StatusCode)
		}

		var response PostsResponse
		if err := json.Unmarshal(bodyBytes, &response); err != nil {
			var apiError struct {
				Error   string `json:"error"`
				Message string `json:"message"`
			}
			if jsonErr := json.Unmarshal(bodyBytes, &apiError); jsonErr == nil {
				return nil, fmt.Errorf("WordPress API error: %s - %s", apiError.Error, apiError.Message)
			}

			logger.Log.Errorf("Error decoding posts response: %v, body: %s", err, string(bodyBytes))
			return nil, fmt.Errorf("failed to decode API response: %v", err)
		}

		if response.Posts == nil {
			continue
		}

		for _, p := range response.Posts {
			if p.ModifiedTime().After(modifiedAfter) {
				allPosts = append(allPosts, p)
			}
		}
	}

	logger.Log.Infof("Total posts fetched: %d", len(allPosts))
	return allPosts, nil
}
