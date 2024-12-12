package wpcom

import "time"

type Post struct {
	ID           int    `json:"ID"`
	Date         string `json:"date"`
	Modified     string `json:"modified"`
	Title        string `json:"title"`
	ContentPlain string `json:"content_plain"`
}

func (p Post) ModifiedTime() time.Time {
	t, _ := time.Parse(time.RFC3339, p.Modified)
	return t
}

type PostsResponse struct {
	Found int    `json:"found"`
	Posts []Post `json:"posts"`
}
