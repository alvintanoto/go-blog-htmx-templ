package entity

import (
	"database/sql"
	"time"
)

const (
	PostVisibilityPrivate   int = 0
	PostVisibilityFollowers int = 1
	PostVisibilityPublic    int = 2
)

type Post struct {
	ID              int
	UserID          string
	Content         string
	ReplyCount      int
	LikeCount       int
	DislikeCount    int
	ImpressionCount int
	SaveCount       int
	Visibility      int
	ReplyTo         *string
	IsDraft         bool
	PostedAt        *time.Time
	CreatedAt       time.Time
	CreatedBy       *string
	UpdatedAt       time.Time
	UpdatedBy       *string
	IsDeleted       bool

	// joins with user name
	Username string
}

func (p *Post) Scan(row *sql.Row) (err error) {
	err = row.Scan(
		&p.ID,
		&p.UserID,
		&p.Content,
		&p.ReplyCount,
		&p.LikeCount,
		&p.DislikeCount,
		&p.ImpressionCount,
		&p.SaveCount,
		&p.Visibility,
		&p.ReplyTo,
		&p.IsDraft,
		&p.PostedAt,
		&p.CreatedAt,
		&p.CreatedBy,
		&p.UpdatedAt,
		&p.UpdatedBy,
		&p.IsDeleted,
	)

	return err
}

func (p *Post) ScanRows(rows *sql.Rows) (err error) {
	err = rows.Scan(
		&p.ID,
		&p.UserID,
		&p.Content,
		&p.ReplyCount,
		&p.LikeCount,
		&p.DislikeCount,
		&p.ImpressionCount,
		&p.SaveCount,
		&p.Visibility,
		&p.ReplyTo,
		&p.IsDraft,
		&p.PostedAt,
		&p.CreatedAt,
		&p.CreatedBy,
		&p.UpdatedAt,
		&p.UpdatedBy,
		&p.IsDeleted,
	)

	return err
}

func (p *Post) ScanJoinUserRows(rows *sql.Rows) (err error) {
	err = rows.Scan(
		&p.ID,
		&p.UserID,
		&p.Content,
		&p.ReplyCount,
		&p.LikeCount,
		&p.DislikeCount,
		&p.ImpressionCount,
		&p.SaveCount,
		&p.Visibility,
		&p.ReplyTo,
		&p.IsDraft,
		&p.PostedAt,
		&p.CreatedAt,
		&p.CreatedBy,
		&p.UpdatedAt,
		&p.UpdatedBy,
		&p.IsDeleted,
		&p.Username,
	)

	return err
}
