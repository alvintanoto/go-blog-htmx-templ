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
	router.HandleFunc("/load-more-posts", a.Controller.ViewController.HomepageInfiniteScrollHandler())

	authRoute := router.PathPrefix("/auth/").Subrouter()
	authRoute.Use(a.Controller.Middlewares.IsAuthenticated)
	{
		authRoute.HandleFunc("/sign-in", a.Controller.ViewController.SignInHandler())
		authRoute.HandleFunc("/register", a.Controller.ViewController.RegisterHandler())
	}

	postRoute := router.PathPrefix("/post/").Subrouter()
	postRoute.Use(a.Controller.Middlewares.IsAuthenticated)
	{
		postRoute.HandleFunc("/new-post", a.Controller.ViewController.CreatePostHandler())
	}

	draftRoute := router.PathPrefix("/draft/").Subrouter()
	draftRoute.Use(a.Controller.Middlewares.IsAuthenticated)
	{
		draftRoute.HandleFunc("/", a.Controller.ViewController.DraftHandler())
	}

	settingsRoute := router.PathPrefix("/settings/").Subrouter()
	settingsRoute.Use(a.Controller.Middlewares.IsAuthenticated)
	{
		settingsRoute.HandleFunc("/", a.Controller.ViewController.SettingsHandler())
	}

	profileRoute := router.PathPrefix("/profile/").Subrouter()
	profileRoute.Use(a.Controller.Middlewares.IsAuthenticated)
	{
		profileRoute.HandleFunc("/", a.Controller.ViewController.ProfileHandler())
		profileRoute.HandleFunc("/load-more-posts", a.Controller.ViewController.ProfilePostInfiniteScrollHandler())
	}

	// apiRoute := router.PathPrefix("/api/").Subrouter()
	// {
	// 	postApiRoute := apiRoute.PathPrefix("/post/").Subrouter()
	// }

	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./view/assets/"))))

	router.NotFoundHandler = http.HandlerFunc(a.Controller.ViewController.NotFoundViewHandler())

	// mux.Handle("*")
	a.Router = router
}
