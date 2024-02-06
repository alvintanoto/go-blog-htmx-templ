package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"alvintanoto.id/blog-htmx-templ/internal/entity"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) CreateNewPost(userID, content string, isDraft bool) (err error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Println("failed to create new uuid: ", err.Error())
		return nil
	}

	post := &entity.Post{
		ID:         uuid.String(),
		UserID:     userID,
		Content:    content,
		Visibility: entity.PostVisibilityPublic,
		CreatedBy:  userID,
		IsDraft:    isDraft,
		IsDeleted:  false,
	}

	args := []interface{}{
		post.ID,
		post.UserID,
		post.Content,
		post.Visibility,
		post.IsDraft,
		post.CreatedBy,
	}

	var query string
	if isDraft {
		query = `INSERT INTO posts(
			id,
			user_id,
			content,
			visibility,
			is_draft,
			created_by
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		)`
	} else {
		post.PostedAt = time.Now()
		args = append(args, post.PostedAt)
		query = `INSERT INTO posts(
			id,
			user_id,
			content,
			visibility,
			is_draft,
			created_by,
			posted_at
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7
		)`
	}

	row := r.db.QueryRow(query, args...)
	if row.Err() != nil {
		switch e := row.Err().(type) {
		case *pq.Error:
			switch e.Code {
			case "23505":
				fmt.Println("error creating new post: ", ErrorConstraintViolation)
				return ErrorConstraintViolation
			case "23503":
				fmt.Println("error creating new post: ", ErrorForeignKeyViolation)
				return ErrorForeignKeyViolation
			}
		}

		log.Println("failed to create new post: ", row.Err())
		return row.Err()
	}

	return nil
}

func (r *PostRepository) GetUserPostByUserID(userID string, page int) (posts []*entity.Post, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var limit int = 15
	var offset = page * limit

	query := `SELECT 
				id, user_id, content, reply_count, like_count, dislike_count, impressions, save_count, posted_at
			FROM 
				posts
			WHERE
				user_id=$1 AND
				reply_to is null AND
				is_draft=false AND
				is_deleted=false
			ORDER BY
				posted_at DESC
			LIMIT 
				$2
			OFFSET
				$3`

	args := []interface{}{
		userID,
		limit,
		offset,
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Println("failed to query user profile post: ", err.Error())
		return nil, err
	}

	for rows.Next() {
		var entity = new(entity.Post)
		err = rows.Scan(&entity.ID, &entity.UserID, &entity.Content, &entity.ReplyCount, &entity.LikeCount, &entity.DislikeCount, &entity.ImpressionCount, &entity.SaveCount, &entity.PostedAt)
		if err != nil {
			log.Println("failed to scan user profile post: ", err.Error())
			return nil, err
		}

		posts = append(posts, entity)
	}

	return posts, nil
}
