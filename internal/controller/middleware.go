package controller

import (
	"log"
	"net/http"
)

type Middlewares struct {
}

func NewMiddleware() *Middlewares {
	return &Middlewares{}
}

func (m *Middlewares) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
