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
		authRoute.HandleFunc("/sign-in", a.Controller.ViewController.SignInHandler())
		authRoute.HandleFunc("/register", a.Controller.ViewController.RegisterHandler())
	}

	postRoute := router.PathPrefix("/post/").Subrouter()
	postRoute.Use(a.Controller.Middlewares.IsAuthenticated)
	{
		postRoute.HandleFunc("/new-post", a.Controller.ViewController.CreateNewPostViewHandler())
		postRoute.HandleFunc("/new-post", a.Controller.ViewController.CreatePostHandler()).Methods(http.MethodPost)
	}

	profileRoute := router.PathPrefix("/profile/").Subrouter()
	profileRoute.Use(a.Controller.Middlewares.IsAuthenticated)
	{
		profileRoute.HandleFunc("/", a.Controller.ViewController.ProfileHandler())
	}

	settingsRoute := router.PathPrefix("/settings/").Subrouter()
	settingsRoute.Use(a.Controller.Middlewares.IsAuthenticated)
	{
		settingsRoute.HandleFunc("/", a.Controller.ViewController.SettingsHandler())
	}

	apiRoute := router.PathPrefix("/api/").Subrouter()
	{
		apiRoute.HandleFunc("/sign-in", a.Controller.ApiController.SignIn()).Methods(http.MethodPost)
		apiRoute.HandleFunc("/register", a.Controller.ApiController.Register()).Methods(http.MethodPost)

		postApiRoute := router.PathPrefix("/post/").Subrouter()
		postApiRoute.HandleFunc("/preview-post", a.Controller.ApiController.PreviewPost()).Methods(http.MethodPost)
	}

	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./view/assets/"))))

	router.NotFoundHandler = http.HandlerFunc(a.Controller.ViewController.NotFoundViewHandler())

	// mux.Handle("*")
	a.Router = router
}
