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
	router.HandleFunc("/load-posts", a.Controller.ViewController.HomepageInfiniteScrollHandler())

	authRoute := router.PathPrefix("/auth/").Subrouter()
	authRoute.Use(a.Controller.Middlewares.IsAuthenticated)
	{
		authRoute.HandleFunc("/sign-in", a.Controller.ViewController.SignInHandler())
		authRoute.HandleFunc("/register", a.Controller.ViewController.RegisterHandler())
	}

	postRoute := router.PathPrefix("/post/").Subrouter()
	postRoute.Use(a.Controller.Middlewares.IsAuthenticated)
	{
		postRoute.HandleFunc("/post_new_post", a.Controller.ViewController.PostNewPostHandler())
		postRoute.HandleFunc("/content", a.Controller.ViewController.PostContentHandler())
		postRoute.HandleFunc("/{id}", a.Controller.ViewController.PostDetailHandler())
	}

	settingsRoute := router.PathPrefix("/settings/").Subrouter()
	settingsRoute.Use(a.Controller.Middlewares.IsAuthenticated)
	{
		settingsRoute.HandleFunc("/", a.Controller.ViewController.SettingsHandler())
		settingsRoute.HandleFunc("/appearance/theme", a.Controller.ViewController.SettingsChangeThemeHandler())
		settingsRoute.HandleFunc("/sign-out-confirmation", a.Controller.ViewController.ShowSignOutConfirmation())
		settingsRoute.HandleFunc("/sign-out", a.Controller.ViewController.SignOutHandler())
	}

	profileRoute := router.PathPrefix("/profile/").Subrouter()
	profileRoute.Use(a.Controller.Middlewares.IsAuthenticated)
	{
		profileRoute.HandleFunc("/", a.Controller.ViewController.ProfileHandler())
		profileRoute.HandleFunc("/load-posts", a.Controller.ViewController.ProfilePostInfiniteScrollHandler())
	}

	// apiRoute := router.PathPrefix("/api/").Subrouter()
	// {
	// 	postApiRoute := apiRoute.PathPrefix("/post/").Subrouter()
	// }

	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./public/assets/"))))

	router.NotFoundHandler = http.HandlerFunc(a.Controller.ViewController.NotFoundViewHandler())

	// mux.Handle("*")
	a.Router = router
}
