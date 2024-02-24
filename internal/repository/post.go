package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"alvintanoto.id/blog-htmx-templ/internal/entity"
	"github.com/lib/pq"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) CreateNewPost(userID, content string, isDraft bool) (err error) {
	if err != nil {
		log.Println("failed to create new uuid: ", err.Error())
		return nil
	}

	currentTime := time.Now()
	post := &entity.Post{
		UserID:     userID,
		Content:    content,
		Visibility: entity.PostVisibilityPublic,
		CreatedBy:  &userID,
		IsDraft:    isDraft,
		IsDeleted:  false,
		PostedAt:   &currentTime,
	}

	args := []interface{}{
		post.UserID,
		post.Content,
		post.Visibility,
		post.IsDraft,
		post.CreatedBy,
		post.PostedAt,
	}

	query := `INSERT INTO posts(
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
			$6
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

func (r *PostRepository) GetPost(postID string) (post *entity.Post, err error) {
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
				reply_to is null AND
				is_deleted=false
			`

	args := []interface{}{
		postID,
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

func (r *PostRepository) GetUserPost(userID, postID string) (post *entity.Post, err error) {
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

func (r *PostRepository) GetPublicTimeline() (posts []*entity.Post, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var limit int = 15

	query := `SELECT 
				id, user_id, content, reply_count, like_count, dislike_count, impressions, 
				save_count, visibility, reply_to, is_draft, posted_at, created_at,
				created_by, updated_at, updated_by, is_deleted
			FROM 
				posts
			WHERE
				reply_to is null AND
				is_draft=false AND
				is_deleted=false AND
				visibility = $1
			ORDER BY
				id DESC
			LIMIT 
				$2
			`

	args := []interface{}{
		entity.PostVisibilityPublic,
		limit,
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

func (r *PostRepository) GetMorePublicTimeline(lastPosition int) (posts []*entity.Post, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var limit int = 15

	query := `SELECT 
				id, user_id, content, reply_count, like_count, dislike_count, impressions, 
				save_count, visibility, reply_to, is_draft, posted_at, created_at,
				created_by, updated_at, updated_by, is_deleted
			FROM 
				posts
			WHERE
				reply_to is null AND
				is_draft=false AND
				is_deleted=false AND
				visibility = $1 AND
				id < $2
			ORDER BY
				id DESC
			LIMIT 
				$3
			`

	args := []interface{}{
		entity.PostVisibilityPublic,
		lastPosition,
		limit,
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

func (r *PostRepository) GetUserPosts(userID string) (posts []*entity.Post, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var limit int = 15

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
				id DESC
			LIMIT 
				$2
			`

	args := []interface{}{
		userID,
		limit,
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

func (r *PostRepository) GetMoreUserPosts(userID string, lastPosition int) (posts []*entity.Post, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var limit int = 15

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
				is_deleted=false AND
				id < $2
			ORDER BY
				id DESC
			LIMIT 
				$3
			`

	args := []interface{}{
		userID,
		lastPosition,
		limit,
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

func (r *PostRepository) CreateBatch(contents []string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	userIds := []string{
		"8c5902e3-9354-47ce-a4eb-d3b339b2205b",
		"1015cff7-f938-4d76-aada-83575e0fcf98",
		"b93ee4d0-600b-48a3-9557-ad605d9e9672",
		"ac396233-fd77-4cfa-a1d7-90afa3c236b6",
		"fcb9e0bb-a85c-4e61-be48-46488d6deb28",
		"595f931f-9db5-41b7-b49f-126382083d2d",
		"52543c55-a982-4d1c-8b10-be4325766a12",
	}

	stmt, _ := tx.Prepare(pq.CopyIn("posts", "user_id", "content", "visibility", "posted_at", "created_at", "created_by")) // MessageDetailRecord is the table name
	for i, content := range contents {
		fmt.Println("inserting:", i)
		inputUserID := userIds[rand.Intn(len(userIds))]

		_, err := stmt.Exec(inputUserID, content, entity.PostVisibilityPublic, time.Now(), time.Now(), inputUserID)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = stmt.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = tx.Commit()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
