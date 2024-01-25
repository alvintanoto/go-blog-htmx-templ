package controller

import (
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
			http.Redirect(w, r, "/register", http.StatusPermanentRedirect)
			return
		}
		err = decoder.Decode(&payload, r.PostForm)
		if err != nil {
			log.Println("error decoding payload: ", err.Error())
			if err != nil {
				log.Println("error saving session: ", err.Error())
				http.Redirect(w, r, "/register", http.StatusPermanentRedirect)
				return
			}
		}

		errCount := 0

		username := strings.TrimSpace(payload.Username)
		email := strings.TrimSpace(payload.Email)
		password := payload.Password
		confirmPassword := payload.ConfirmPassword

		if username == "" {
			session.AddFlash("Username must not be empty.", "username")
			errCount += 1
		}

		if email == "" {
			session.AddFlash("Email must not be empty.", "email")
			errCount += 1
		}

		if password == "" {
			session.AddFlash("Password must not be empty.", "password")
			errCount += 1
		}

		if confirmPassword == "" {
			session.AddFlash("Confirm password must not be empty.", "confirm_password")
			errCount += 1
		}

		if len(username) <= 6 {
			session.AddFlash("Username length must be more than 6 character.", "username")
			errCount += 1
		}

		if len(email) <= 10 {
			session.AddFlash("Email length must be more than 10 character.", "email")
			errCount += 1
		}

		if len(password) <= 6 {
			session.AddFlash("Password length must be more than 6 character.", "password")
			errCount += 1
		}

		if password != confirmPassword {
			session.AddFlash("Password and confirm password mismatch.", "password")
			session.AddFlash("Password and confirm password mismatch.", "confirm_password")
			errCount += 1
		}

		if errCount > 0 {
			err := sessions.Save(r, w)
			if err != nil {
				log.Println("err saving session:", err.Error())
			}
			http.Redirect(w, r, "/register", http.StatusPermanentRedirect)
			return
		}

		user, err := ac.Service.UserService.RegisterUser(username, email, password)
		if err != nil {
			log.Println("failed registering new user: ", err.Error())
			session.AddFlash("Failed registering new user, please try again later.", "error")
			sessions.Save(r, w)
			http.Redirect(w, r, "/register", http.StatusPermanentRedirect)
			return
		}

		session.Values["user"] = &dto.UserDTO{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		}
		sessions.Save(r, w)
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}
}
