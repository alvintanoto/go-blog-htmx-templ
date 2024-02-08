package dto

type PostDTO struct {
	ID          string    `json:"id"`
	Content     string    `json:"content"`
	Replies     []PostDTO `json:"replies"`
	ReplyCounts int       `json:"reply_counts"`
	Likes       int       `json:"likes"`
	Dislikes    int       `json:"dislikes"`
	Impressions int       `json:"impressions"`
	SavedCounts int       `json:"saved_count"`
	Poster      UserDTO   `json:"poster"`
	PostedAt    string    `json:"posted_at"`
	CreatedAt   string    `json:"created_at"`
}

type SubmitPostDTO struct {
	Content    string `schema:"content"`
	SubmitType string `schema:"submit_type"`
}
