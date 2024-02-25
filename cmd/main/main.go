package main

import (
	"context"
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"alvintanoto.id/blog-htmx-templ/internal/controller"
	"alvintanoto.id/blog-htmx-templ/internal/db"
	"alvintanoto.id/blog-htmx-templ/internal/dto"
	"alvintanoto.id/blog-htmx-templ/internal/repository"
	"alvintanoto.id/blog-htmx-templ/internal/service"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/v9"
)

type Application struct {
	Configurations *Configurations
	Database       *sql.DB
	Store          *redisstore.RedisStore

	Router     *mux.Router
	Controller *controller.Controller
	Service    *service.Service
	Repository *repository.Repository
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
	a.Controller = controller.NewController(a.Store, a.Service)
}

func (a *Application) InitializeService() {
	a.Service = service.NewService(a.Repository)
}

func (a *Application) InitializeRepository() {
	a.Repository = repository.NewRepository(a.Database)
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

func (a *Application) InitializeSession() {
	client := redis.NewClient(&redis.Options{
		Addr: a.Configurations.Redis.Url,
	})

	// New default RedisStore
	store, err := redisstore.NewRedisStore(context.Background(), client)
	if err != nil {
		log.Fatal("failed to create redis store: ", err)
	}

	store.Options(sessions.Options{
		Path:   "/",
		MaxAge: 60 * 60 * 3,
	})

	a.Store = store

	// register structs for sessions
	gob.Register(&dto.UserDTO{})
}

func main() {
	app := &Application{}
	app.InitializeConfigs()
	app.InitializeDatabase()
	app.InitializeSession()
	app.InitializeRepository()
	app.InitializeService()
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
