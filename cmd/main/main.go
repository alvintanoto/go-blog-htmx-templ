package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"alvintanoto.id/blog-htmx-templ/internal/db"
	"alvintanoto.id/blog-htmx-templ/view"
	"github.com/gorilla/mux"
)

type Application struct {
	Configurations *Configurations
	Database       *sql.DB
	Router         *mux.Router
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

// InitializeDatabase to set up connection to database
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

// SetupRoutes to setup routes here
func (a *Application) SetupRoutes() {
	router := mux.NewRouter()
	router.HandleFunc("/", http.HandlerFunc((func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("homepage"))
	})))

	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		view.NotFoundPage().Render(r.Context(), w)
	})

	// mux.Handle("*")
	a.Router = router
}

func main() {
	app := &Application{}
	app.InitializeConfigs()
	app.InitializeDatabase()
	app.SetupRoutes()

	// start the server here
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", app.Configurations.Server.Port),
		Handler:      app.Router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Println("starting server on", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to start server: #%v", err)
	}
}
