package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

type Middlewares struct {
	Store *sessions.CookieStore
}

func NewMiddleware(store *sessions.CookieStore) *Middlewares {
	return &Middlewares{
		Store: store,
	}
}

func (m *Middlewares) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func (m *Middlewares) IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("path", r.URL.Path)
		store, _ := m.Store.Get(r, "default")
		user := store.Values["user"]

		if user != nil && (r.URL.Path == "/sign-in" || r.URL.Path == "/register") {
			http.Redirect(w, r, "/", http.StatusPermanentRedirect)
			return
		}

		next.ServeHTTP(w, r)
	})
}
