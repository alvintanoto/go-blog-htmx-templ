package repository

import (
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
