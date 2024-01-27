package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"alvintanoto.id/blog-htmx-templ/internal/entity"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (user *entity.User, err error) {
	user = new(entity.User)
	query := `SELECT id, username, email, password FROM blog_user WHERE username=$1 AND is_deleted=false`
	args := []interface{}{
		username,
	}

	row := r.db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		log.Println("error querying user by username: ", row.Err())
		return nil, err
	}

	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		log.Println("fail scanning row to struct: ", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorRecordNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) RegisterUser(username, email, password string) (user *entity.User, err error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Println("failed to create new uuid: ", err.Error())
		return nil, err
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("failed to create password bcrypt: ", err.Error())
		return nil, err
	}

	user = &entity.User{
		ID:           uuid.String(),
		Username:     username,
		Email:        email,
		Password:     string(encryptedPassword),
		CreatedBy:    username,
		UpdatedBy:    username,
		LastLoggedIn: time.Now(),
		IsDeleted:    false,
	}

	query := `INSERT INTO blog_user(
		id,
		username,
		email,
		password,
		created_by,
		updated_by,
		last_logged_in,
		is_deleted
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8
	)`

	args := []interface{}{
		user.ID,
		user.Username,
		user.Email,
		user.Password,
		user.CreatedBy,
		user.UpdatedBy,
		user.LastLoggedIn,
		user.IsDeleted,
	}

	row := r.db.QueryRow(query, args...)
	if row.Err() != nil {
		switch e := row.Err().(type) {
		case *pq.Error:
			switch e.Code {
			case "23505":
				fmt.Println("error creating new user: ", ErrorConstraintViolation)
				return nil, ErrorConstraintViolation
			case "23503":
				fmt.Println("error creating new user: ", ErrorForeignKeyViolation)
				return nil, ErrorForeignKeyViolation
			}
		}

		log.Println("failed to create user: ", row.Err())
		return nil, row.Err()
	}

	return user, nil
}
