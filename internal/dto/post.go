package dto

import "time"

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
	PostedAt    time.Time `json:"posted_at"`
}

type PreviewPostDTO struct {
	Content string `json:"content"`
	Preview bool   `json:"preview"`
}
