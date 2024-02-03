package controller

import (
	"fmt"
	"log"
	"net/http"

	"alvintanoto.id/blog-htmx-templ/internal/dto"
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
		store, _ := m.Store.Get(r, "default")
		user := store.Values["user"]

		if user != nil {
			userDTO := user.(*dto.UserDTO)
			log.Println(fmt.Sprintf("[%s - %s]", userDTO.ID, r.RemoteAddr), r.Method, r.RequestURI)
		} else {
			log.Println(fmt.Sprintf("[%s]", r.RemoteAddr), r.Method, r.RequestURI)
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middlewares) IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		store, _ := m.Store.Get(r, "default")
		user := store.Values["user"]

		if user != nil {
			if r.URL.Path == "/auth/sign-in" || r.URL.Path == "/auth/register" {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
			return
		}

		// user exist
		if r.URL.Path != "/auth/sign-in" && r.URL.Path != "/auth/register" && r.URL.Path != "/" {
			http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
