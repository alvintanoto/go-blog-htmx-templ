package dto

import "time"

type PostDTO struct {
	ID          string    `json:"id"`
	Message     string    `json:"post"`
	Replies     []PostDTO `json:"replies"`
	ReplyCounts int       `json:"reply_counts"`
	Likes       int       `json:"likes"`
	Impressions int       `json:"impressions"`
	SavedCounts int       `json:"saved_count"`
	Poster      UserDTO   `json:"poster"`
	PostedAt    time.Time `json:"posted_at"`
}
