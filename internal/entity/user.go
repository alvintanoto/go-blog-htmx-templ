package entity

import (
	"database/sql"
	"time"
)

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

type UserConfig struct {
	Key   string
	Value string
}

func (uc *UserConfig) ScanUserConfig(rows *sql.Rows) (err error) {
	err = rows.Scan(
		&uc.Key,
		&uc.Value,
	)

	return err
}
