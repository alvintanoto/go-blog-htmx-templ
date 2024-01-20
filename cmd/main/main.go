package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"alvintanoto.id/blog-htmx-templ/internal/db"
)

type Application struct {
	Configurations *Configurations
	Database       *sql.DB
}

// InitializeConfigs set up env variable configurations
func (a *Application) InitializeConfigs() {
	configs := &Configurations{}
	err := configs.ReadConfigurations()
	if err != nil {
		log.Println("failed to read configuration: ", err.Error())
		return
	}

	a.Configurations = configs
}

// Initialize Database to set up connection to database
func (a *Application) InitializeDatabase() {
	db, err := db.InitializeDB(a.Configurations.Database.User,
		a.Configurations.Database.Password,
		a.Configurations.Database.Host,
		a.Configurations.Database.Port,
		a.Configurations.Database.Name)

	if err != nil {
		log.Println("failed to initialize database: ", err.Error())
		return
	}

	a.Database = db
}

func main() {
	app := &Application{}
	app.InitializeConfigs()
	app.InitializeDatabase()

	// start the server here
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", app.Configurations.Server.Port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Println("starting server on", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to start server: #%v", err)
	}
}
