package entity

import "time"

const (
	PostVisibilityPrivate   int = 0
	PostVisibilityFollowers int = 1
	PostVisibilityPublic    int = 2
)

type Post struct {
	ID              string
	UserID          string
	Content         string
	ReplyCount      int
	LikeCount       int
	DislikeCount    int
	ImpressionCount int
	SaveCount       int
	Visibility      int
	ReplyTo         string
	IsDraft         bool
	PostedAt        time.Time
	CreatedAt       time.Time
	CreatedBy       string
	UpdatedAt       time.Time
	UpdatedBy       string
	IsDeleted       bool
}
