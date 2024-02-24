package controller

import (
	"alvintanoto.id/blog-htmx-templ/internal/service"
	"github.com/rbcervilla/redisstore/v9"
)

type Controller struct {
	Middlewares *Middlewares

	ViewController ViewController
	ApiController  *ApiController
}

func NewController(store *redisstore.RedisStore, service *service.Service) *Controller {
	return &Controller{
		Middlewares:    NewMiddleware(store),
		ViewController: NewViewController(service, store),
		ApiController:  NewApiController(service, store),
	}
}
