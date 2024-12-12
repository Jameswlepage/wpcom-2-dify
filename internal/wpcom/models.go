package wpcom

import (
	"dify-wp-sync/internal/logger"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

// Post represents a WordPress.com post.
type Post struct {
	ID       int    `json:"ID"`
	Date     string `json:"date"`
	Modified string `json:"modified"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

func (p Post) ModifiedTime() time.Time {
	t, _ := time.Parse(time.RFC3339, p.Modified)
	return t
}

func (p Post) GetMarkdownContent() string {
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(p.Content)
	if err != nil {
		// Log error but return original content as fallback
		logger.Log.Errorf("Failed to convert HTML to Markdown for post %d: %v", p.ID, err)
		return p.Content
	}
	return markdown
}

// PostsResponse represents the WordPress.com API response for posts.
type PostsResponse struct {
	Found int          `json:"found"`
	Posts []Post       `json:"posts"`
	Meta  PostMetadata `json:"meta"`
}

// PostMetadata represents the metadata returned in the WordPress.com API response
type PostMetadata struct {
	Links struct {
		Counts string `json:"counts"`
	} `json:"links"`
	WPCom bool `json:"wpcom"`
}
