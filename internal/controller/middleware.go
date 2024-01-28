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
		store, _ := m.Store.Get(r, "default")
		user := store.Values["user"]

		fmt.Println(user)

		if user != nil {
			if r.URL.Path == "/sign-in" || r.URL.Path == "/register" {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		}

		// user nil
		if r.URL.Path != "/sign-in" && r.URL.Path != "/register" && r.URL.Path != "/" {
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
