package entity

import "time"

type User struct {
	ID           string
	Username     string
	Email        string
	Password     string
	LastLoggedIn time.Time
	CreatedAt    time.Time
	CreatedBy    string
	UpdatedAt    time.Time
	UpdatedBy    string
	IsDeleted    bool
}
