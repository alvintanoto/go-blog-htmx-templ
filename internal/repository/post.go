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

	currentTime := time.Now()
	post := &entity.Post{
		ID:         uuid.String(),
		UserID:     userID,
		Content:    content,
		Visibility: entity.PostVisibilityPublic,
		CreatedBy:  &userID,
		IsDraft:    isDraft,
		IsDeleted:  false,
		PostedAt:   &currentTime,
	}

	args := []interface{}{
		post.ID,
		post.UserID,
		post.Content,
		post.Visibility,
		post.IsDraft,
		post.CreatedBy,
		post.PostedAt,
	}

	query := `INSERT INTO posts(
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

func (r *PostRepository) GetPost(userID, postID string) (post *entity.Post, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	post = new(entity.Post)

	query := `SELECT 
				id, user_id, content, reply_count, like_count, dislike_count, impressions, 
				save_count, visibility, reply_to, is_draft, posted_at, created_at,
				created_by, updated_at, updated_by, is_deleted
			FROM 
				posts
			WHERE
				id=$1 AND 
				user_id=$2 AND
				reply_to is null AND
				is_deleted=false
			`

	args := []interface{}{
		postID,
		userID,
	}

	result := r.db.QueryRowContext(ctx, query, args...)
	if result.Err() != nil {
		log.Println("error querying user post by user id: ", result.Err().Error())
		return nil, result.Err()
	}

	err = post.Scan(result)
	if err != nil {
		log.Println("error scanning user post by user id: ", err.Error())
		return nil, err
	}

	return post, nil
}

func (r *PostRepository) GetPosts(userID string, page int) (posts []*entity.Post, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var limit int = 15
	var offset = page * limit

	query := `SELECT 
				id, user_id, content, reply_count, like_count, dislike_count, impressions, 
				save_count, visibility, reply_to, is_draft, posted_at, created_at,
				created_by, updated_at, updated_by, is_deleted
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
		err = entity.ScanRows(rows)
		if err != nil {
			log.Println("failed to scan user profile post: ", err.Error())
			return nil, err
		}

		posts = append(posts, entity)
	}

	return posts, nil
}

func (r *PostRepository) GetDrafts(userID string, page int) (posts []*entity.Post, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var limit int = 15
	var offset = page * limit

	query := `
			SELECT 
				id, user_id, content, reply_count, like_count, dislike_count, impressions, 
				save_count, visibility, reply_to, is_draft, posted_at, created_at,
				created_by, updated_at, updated_by, is_deleted
			FROM 
				posts
			WHERE
				user_id=$1 AND
				reply_to is null AND
				is_draft=true AND
				is_deleted=false
			ORDER BY
				updated_at DESC
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
		log.Println("failed to query user drafted post: ", err.Error())
		return nil, err
	}

	for rows.Next() {
		var entity = new(entity.Post)
		entity.ScanRows(rows)
		if err != nil {
			log.Println("failed to scan user drafted post: ", err.Error())
			return nil, err
		}

		posts = append(posts, entity)
	}

	return posts, nil
}

func (r *PostRepository) UpdatePost(post *entity.Post) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// currentTime := time.Now()
	// post.Content = payloadContent
	// post.UpdatedAt = currentTime
	// post.UpdatedBy = &userID

	// if !isDrafting {
	// 	post.PostedAt = &currentTime
	// }

	query := `
		UPDATE
			posts
		SET 
			content = $1,
			is_draft = $2,
			updated_at = $3,
			updated_by = $4,
			posted_at = $5
		WHERE 
			id = $6
	`

	args := []interface{}{
		post.Content,
		post.IsDraft,
		post.UpdatedAt,
		post.UpdatedBy,
		post.PostedAt,
		post.ID,
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
