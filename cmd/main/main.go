package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"alvintanoto.id/blog-htmx-templ/internal/controller"
	"alvintanoto.id/blog-htmx-templ/internal/db"
	"github.com/gorilla/mux"
)

type Application struct {
	Configurations *Configurations
	Database       *sql.DB

	Router     *mux.Router
	Controller *controller.Controller
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

func (a *Application) InitializeController() {
	a.Controller = &controller.Controller{}
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

	router.Use(a.Controller.Middlewares.LoggingMiddleware)

	router.HandleFunc("/", a.Controller.ViewController.HomepageViewHandler())
	router.HandleFunc("/sign-in", a.Controller.ViewController.SignInHandler())

	postRoute := router.PathPrefix("/post/").Subrouter()
	{
		postRoute.HandleFunc("/new_post", a.Controller.ViewController.CreateNewPostHandler())
	}

	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	router.NotFoundHandler = http.HandlerFunc(a.Controller.ViewController.NotFoundViewHandler())

	// mux.Handle("*")
	a.Router = router
}

func main() {
	app := &Application{}
	app.InitializeConfigs()
	app.InitializeDatabase()
	app.InitializeController()
	app.SetupRoutes()

	server := http.Server{
		Addr:         fmt.Sprintf(":%d", app.Configurations.Server.Port),
		Handler:      app.Router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// start the server here with goroutine so it doesnt block
	go func() {
		log.Println("starting server on", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	server.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("server shutting down")
	os.Exit(0)
}
