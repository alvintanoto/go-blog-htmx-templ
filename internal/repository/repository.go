package repository

import (
	"database/sql"
	"errors"
)

var (
	ErrorRecordNotFound      = errors.New("record not found")
	ErrorConstraintViolation = errors.New("constraint_violation")
	ErrorForeignKeyViolation = errors.New("foreign_key_violation")
)

type Repository struct {
	UserRepository *UserRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		UserRepository: NewUserRepository(db),
	}
}
