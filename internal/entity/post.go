package entity

import "time"

type Post struct {
	ID                string
	Content           string
	ReplyCount        int
	LikeCount         int
	DislikeCount      int
	ImpressionCount   int
	SaveCount         int
	Visibility        int
	ReplyTo           string
	PreviousVersionID string
	CreatedAt         time.Time
	CreatedBy         string
	UpdatedAt         time.Time
	UpdatedBy         string
	IsDeleted         bool
}
