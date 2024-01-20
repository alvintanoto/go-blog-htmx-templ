package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func InitializeDB(dbUsername, dbPassword, dbHost, dbPort, dbName string) (db *sql.DB, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUsername, dbPassword, dbHost, dbPort, dbName)
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Println("failed to create connection to database: ", err)
		return nil, err
	}

	db.SetConnMaxIdleTime(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	log.Println("Database initialized")
	return db, nil
}
