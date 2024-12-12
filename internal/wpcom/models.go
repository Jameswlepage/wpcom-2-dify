package wpcom

import "time"

// Post represents a WordPress.com post.
// We now use 'content_raw' instead of 'content_plain' for a plain-text like field.
type Post struct {
	ID         int    `json:"ID"`
	Date       string `json:"date"`
	Modified   string `json:"modified"`
	Title      string `json:"title"`
	ContentRaw string `json:"content_raw"`
}

func (p Post) ModifiedTime() time.Time {
	t, _ := time.Parse(time.RFC3339, p.Modified)
	return t
}

type PostsResponse struct {
	Found int    `json:"found"`
	Posts []Post `json:"posts"`
}
