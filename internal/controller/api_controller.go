package controller

import (
	"alvintanoto.id/blog-htmx-templ/internal/service"
	"github.com/rbcervilla/redisstore/v9"
)

type ApiController struct {
	Store   *redisstore.RedisStore
	Service *service.Service
}

func NewApiController(service *service.Service, store *redisstore.RedisStore) *ApiController {
	return &ApiController{
		Service: service,
		Store:   store,
	}
}
