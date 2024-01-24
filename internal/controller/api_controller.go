package controller

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"alvintanoto.id/blog-htmx-templ/internal/dto"
	"alvintanoto.id/blog-htmx-templ/internal/service"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
)

type ApiController struct {
	Store   *sessions.CookieStore
	Service *service.Service
}

func NewApiController(service *service.Service, store *sessions.CookieStore) *ApiController {
	return &ApiController{
		Service: service,
		Store:   store,
	}
}
func (ac *ApiController) Register() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := schema.NewDecoder()
		session, _ := ac.Store.Get(r, "default")

		var payload dto.RegisterUserRequestDTO

		err := r.ParseForm()
		if err != nil {
			log.Println("error parsing form: ", err.Error())
			http.Redirect(w, r, "/register", http.StatusTemporaryRedirect)
			return
		}
		err = decoder.Decode(&payload, r.PostForm)
		if err != nil {
			log.Println("error decoding payload: ", err.Error())
			if err != nil {
				log.Println("error saving session: ", err.Error())
				http.Redirect(w, r, "/register", http.StatusTemporaryRedirect)
				return
			}
		}

		errCount := 0

		if strings.TrimSpace(payload.Username) == "" {
			session.AddFlash("Username must not be empty.", "username")
			errCount += 1
		}

		if strings.TrimSpace(payload.Email) == "" {
			session.AddFlash("Email must not be empty.", "email")
			errCount += 1
		}

		if strings.TrimSpace(payload.Password) == "" {
			session.AddFlash("Password must not be empty.", "password")
			errCount += 1
		}

		if strings.TrimSpace(payload.ConfirmPassword) == "" {
			session.AddFlash("Confirm password must not be empty.", "confirm_password")
			errCount += 1
		}

		if len(strings.TrimSpace(payload.Username)) <= 6 {
			session.AddFlash("Username length must be more than 6 character.", "username")
			errCount += 1
		}

		if len(strings.TrimSpace(payload.Email)) <= 10 {
			session.AddFlash("Email length must be more than 10 character.", "email")
			errCount += 1
		}

		if len(strings.TrimSpace(payload.Password)) <= 6 {
			session.AddFlash("Password length must be more than 6 character.", "password")
			errCount += 1
		}

		if payload.Password != payload.ConfirmPassword {
			session.AddFlash("Password and confirm password mismatch.", "password")
			session.AddFlash("Password and confirm password mismatch.", "confirm_password")
			errCount += 1
		}

		if errCount > 0 {
			err := sessions.Save(r, w)
			if err != nil {
				log.Println("err saving session:", err.Error())
			}
			http.Redirect(w, r, "/register", http.StatusTemporaryRedirect)
			return
		}

		fmt.Println("username", payload.Username)
		fmt.Println("email", payload.Email)
		fmt.Println("password", payload.Password)

		fmt.Println("register")

	}
}
