package wpcom

import (
	"dify-wp-sync/internal/logger"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// WPClient interacts with the WordPress.com API.
type WPClient struct {
	AccessToken string
	SiteID      string
	httpClient  *http.Client
}

// NewWPClient creates a new WPClient with the provided token, site ID, and a default timeout.
func NewWPClient(token, siteID string) *WPClient {
	return &WPClient{
		AccessToken: token,
		SiteID:      siteID,
		httpClient:  &http.Client{Timeout: 30 * time.Second},
	}
}

// GetPosts fetches posts of specified types updated since modifiedAfter.
// Now it includes pagination to fetch all posts/pages, not just the first 100.
func (c *WPClient) GetPosts(modifiedAfter time.Time, postTypes []string) ([]Post, error) {
	var allPosts []Post

	for _, postType := range postTypes {
		logger.Log.Infof("Fetching posts of type '%s' from site %s", postType, c.SiteID)

		// We'll fetch in batches of 100 using an offset
		offset := 0
		limit := 100

		for {
			apiURL := fmt.Sprintf("https://public-api.wordpress.com/rest/v1.1/sites/%s/posts", c.SiteID)
			params := url.Values{}
			params.Set("number", strconv.Itoa(limit))  // how many items to fetch per request
			params.Set("offset", strconv.Itoa(offset)) // how many items to skip
			params.Set("order_by", "modified")
			params.Set("order", "DESC")
			params.Set("fields", "ID,date,modified,title,content,type")
			params.Set("type", postType)

			req, err := http.NewRequest("GET", apiURL+"?"+params.Encode(), nil)
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

			if resp.StatusCode != http.StatusOK {
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

			// Add logging for pagination details
			logger.Log.Infof("Batch stats - Found: %d, Offset: %d, Limit: %d, Posts in response: %d",
				response.Found, offset, limit, len(response.Posts))

			// If no posts returned in this batch, we're done for this post type.
			if len(response.Posts) == 0 {
				logger.Log.Infof("No more posts found for type '%s'", postType)
				break
			}

			// Log posts that match our modified time criteria
			matchingPosts := 0
			for _, p := range response.Posts {
				if p.ModifiedTime().After(modifiedAfter) {
					allPosts = append(allPosts, p)
					matchingPosts++
				}
			}
			logger.Log.Infof("Found %d posts modified after %v in this batch", matchingPosts, modifiedAfter)

			// If offset + limit >= total found, we've fetched all pages/posts for this type.
			if offset+limit >= response.Found {
				logger.Log.Infof("Reached end of posts for type '%s' (offset %d >= total %d)",
					postType, offset+limit, response.Found)
				break
			}
			offset += limit
		}
	}

	logger.Log.Infof("Total posts fetched: %d", len(allPosts))
	return allPosts, nil
}

func (c *WPClient) GetPostsBatch(modifiedAfter time.Time, postType string, offset, limit int) ([]Post, bool, error) {
	logger.Log.Infof("Fetching batch of type '%s' from site %s (offset: %d, limit: %d)",
		postType, c.SiteID, offset, limit)

	apiURL := fmt.Sprintf("https://public-api.wordpress.com/rest/v1.1/sites/%s/posts", c.SiteID)
	params := url.Values{}
	params.Set("number", strconv.Itoa(limit))
	params.Set("offset", strconv.Itoa(offset))
	params.Set("order_by", "modified")
	params.Set("order", "DESC")
	params.Set("fields", "ID,date,modified,title,content,type")
	params.Set("type", postType)

	req, err := http.NewRequest("GET", apiURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("User-Agent", "Dify-WP-Sync/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	logger.Log.Infof("Received response status: %d from WordPress API", resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, false, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("unexpected status code %d from WordPress API", resp.StatusCode)
	}

	var response PostsResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, false, fmt.Errorf("failed to decode API response: %v", err)
	}

	logger.Log.Infof("Batch stats - Found: %d, Offset: %d, Limit: %d, Posts in response: %d",
		response.Found, offset, limit, len(response.Posts))

	var matchingPosts []Post
	for _, p := range response.Posts {
		if p.ModifiedTime().After(modifiedAfter) {
			matchingPosts = append(matchingPosts, p)
		}
	}

	hasMore := offset+limit < response.Found
	return matchingPosts, hasMore, nil
}
