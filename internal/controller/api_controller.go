package controller

import (
	"log"
	"net/http"
	"strings"

	"alvintanoto.id/blog-htmx-templ/internal/dto"
	"alvintanoto.id/blog-htmx-templ/internal/repository"
	"alvintanoto.id/blog-htmx-templ/internal/service"
	"alvintanoto.id/blog-htmx-templ/view"
	"github.com/gomarkdown/markdown"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/crypto/bcrypt"
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

func (ac *ApiController) SignIn() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := ac.Store.Get(r, "default")
		decoder := schema.NewDecoder()

		var payload dto.UserSignInRequestDTO

		err := r.ParseForm()
		if err != nil {
			log.Println("error parsing form:", err.Error())
			http.Redirect(w, r, "/auth/sign-in", http.StatusPermanentRedirect)
			return
		}

		err = decoder.Decode(&payload, r.PostForm)
		if err != nil {
			log.Println("error decoding payload: ", err.Error())
			http.Redirect(w, r, "/auth/sign-in", http.StatusPermanentRedirect)
			return
		}

		if strings.TrimSpace(payload.Username) == "" || strings.TrimSpace(payload.Password) == "" {
			session.AddFlash("Username or password invalid, please try again.", "error")
			sessions.Save(r, w)
			http.Redirect(w, r, "/auth/sign-in", http.StatusPermanentRedirect)
			return
		}

		user, err := ac.Service.AuthenticationService.SignIn(r.Context(), payload.Username, payload.Password)
		if err != nil {
			switch err {
			case repository.ErrorRecordNotFound,
				bcrypt.ErrMismatchedHashAndPassword:
				session.AddFlash("Username or password invalid, please try again.", "error")
			default:
				session.AddFlash("Failed to sign in, please try again later.", "error")
			}
			// TODO: redirect to sign in with flash error
			sessions.Save(r, w)
			http.Redirect(w, r, "/auth/sign-in", http.StatusPermanentRedirect)
			return
		}

		session.Values["user"] = &dto.UserDTO{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		}
		sessions.Save(r, w)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (ac *ApiController) Register() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := ac.Store.Get(r, "default")
		decoder := schema.NewDecoder()

		var payload dto.RegisterUserRequestDTO

		err := r.ParseForm()
		if err != nil {
			log.Println("error parsing form: ", err.Error())
			http.Redirect(w, r, "/auth/register", http.StatusPermanentRedirect)
			return
		}
		err = decoder.Decode(&payload, r.PostForm)
		if err != nil {
			log.Println("error decoding payload: ", err.Error())
			http.Redirect(w, r, "/auth/register", http.StatusPermanentRedirect)
			return
		}

		// TODO: this is temporary solution fix it later
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
			http.Redirect(w, r, "/auth/register", http.StatusPermanentRedirect)
			return
		}

		user, err := ac.Service.UserService.RegisterUser(username, email, password)
		if err != nil {
			log.Println("failed registering new user: ", err.Error())
			switch err {
			case repository.ErrorConstraintViolation:
				session.AddFlash("Username already used, please try another username.", "error")
			default:
				session.AddFlash("Failed registering new user, please try again later.", "error")
			}
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

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (ac *ApiController) PreviewPost() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := schema.NewDecoder()

		var payload dto.PreviewPostDTO
		err := r.ParseForm()
		if err != nil {
			log.Println("error parsing form: ", err.Error())
			http.Redirect(w, r, "/post/new-post", http.StatusPermanentRedirect)
			return
		}

		err = decoder.Decode(&payload, r.PostForm)
		if err != nil {
			log.Println("error decoding payload: ", err.Error())
			http.Redirect(w, r, "/post/new-post", http.StatusPermanentRedirect)
			return
		}

		if payload.Preview {
			maybeUnsafeHTML := markdown.ToHTML([]byte(payload.Value), nil, nil)
			html := bluemonday.UGCPolicy().SanitizeBytes(maybeUnsafeHTML)

			view.PreviewPostContainer(string(html), payload.Value).Render(r.Context(), w)
			return
		}

		view.EditorContainer(payload.Value).Render(r.Context(), w)
	}
}
