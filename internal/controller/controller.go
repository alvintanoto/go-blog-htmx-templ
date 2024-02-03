package controller

import (
	"alvintanoto.id/blog-htmx-templ/internal/service"
	"github.com/gorilla/sessions"
)

type Controller struct {
	Middlewares *Middlewares

	ViewController *ViewController
	ApiController  *ApiController
}

func NewController(store *sessions.CookieStore, service *service.Service) *Controller {
	return &Controller{
		Middlewares:    NewMiddleware(store),
		ViewController: NewViewController(service, store),
		ApiController:  NewApiController(service, store),
	}
}
