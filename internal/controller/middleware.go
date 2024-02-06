package controller

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"alvintanoto.id/blog-htmx-templ/internal/dto"
	"github.com/rbcervilla/redisstore/v9"
)

type Middlewares struct {
	Store *redisstore.RedisStore
}

func NewMiddleware(store *redisstore.RedisStore) *Middlewares {
	return &Middlewares{
		Store: store,
	}
}

func (m *Middlewares) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.RequestURI, "/assets/") {
			next.ServeHTTP(w, r)
			return
		}

		store, _ := m.Store.Get(r, "default")
		userStore := store.Values["user"]

		if userStore != nil {
			user := userStore.(*dto.UserDTO)

			go func() {
				log.Println("===============================")
				log.Println("[user_id]\t:", user.ID)
				log.Println("[addr]\t:", r.RemoteAddr)
				log.Println("[agent]\t:", r.UserAgent())
				log.Println("[method]\t:", r.Method)
				log.Println("[uri]\t:", r.RequestURI)
				log.Println("===============================")
			}()

		} else {
			go func() {
				log.Println("===============================")
				log.Println("[user_id]\t:", "non-registered user")
				log.Println("[addr]\t:", r.RemoteAddr)
				log.Println("[agent]\t:", r.UserAgent())
				log.Println("[method]\t:", r.Method)
				log.Println("[uri]\t:", r.RequestURI)
				log.Println("===============================")
			}()
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middlewares) IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		store, _ := m.Store.Get(r, "default")
		user := store.Values["user"]
		fmt.Println("isauth", user)

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
