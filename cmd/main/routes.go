package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// SetupRoutes to setup routes here
func (a *Application) SetupRoutes() {
	router := mux.NewRouter()

	router.Use(a.Controller.Middlewares.LoggingMiddleware)

	router.HandleFunc("/", a.Controller.ViewController.HomepageViewHandler())

	authRoute := router.PathPrefix("/auth/").Subrouter()
	authRoute.Use(a.Controller.Middlewares.IsAuthenticated)
	{
		authRoute.HandleFunc("/sign-in", a.Controller.ViewController.SignInPostHandler()).Methods(http.MethodPost)
		authRoute.HandleFunc("/sign-in", a.Controller.ViewController.SignInHandler()).Methods(http.MethodGet)
		authRoute.HandleFunc("/register", a.Controller.ViewController.RegisterHandler()).Methods(http.MethodGet)
	}

	postRoute := router.PathPrefix("/post/").Subrouter()
	postRoute.Use(a.Controller.Middlewares.IsAuthenticated)
	{
		postRoute.HandleFunc("/new-post", a.Controller.ViewController.CreatePostHandler()).Methods(http.MethodPost)
		postRoute.HandleFunc("/new-post", a.Controller.ViewController.CreateNewPostViewHandler()).Methods(http.MethodGet)
	}

	profileRoute := router.PathPrefix("/profile/").Subrouter()
	profileRoute.Use(a.Controller.Middlewares.IsAuthenticated)
	{
		profileRoute.HandleFunc("/", a.Controller.ViewController.ProfileHandler()).Methods(http.MethodGet)
	}

	settingsRoute := router.PathPrefix("/settings/").Subrouter()
	settingsRoute.Use(a.Controller.Middlewares.IsAuthenticated)
	{
		settingsRoute.HandleFunc("/", a.Controller.ViewController.SettingsHandler()).Methods(http.MethodGet)
	}

	apiRoute := router.PathPrefix("/api/").Subrouter()
	{
		apiRoute.HandleFunc("/register", a.Controller.ApiController.Register()).Methods(http.MethodPost)

		postApiRoute := router.PathPrefix("/post/").Subrouter()
		postApiRoute.HandleFunc("/preview-post", a.Controller.ApiController.PreviewPost()).Methods(http.MethodPost)
	}

	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./view/assets/"))))

	router.NotFoundHandler = http.HandlerFunc(a.Controller.ViewController.NotFoundViewHandler())

	// mux.Handle("*")
	a.Router = router
}
