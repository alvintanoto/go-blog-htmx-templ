package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"alvintanoto.id/blog-htmx-templ/internal/dto"
	"alvintanoto.id/blog-htmx-templ/internal/repository"
	"alvintanoto.id/blog-htmx-templ/internal/service"
	"alvintanoto.id/blog-htmx-templ/view"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type ViewController struct {
	Session *sessions.CookieStore
	Service *service.Service
}

func NewViewController(service *service.Service, session *sessions.CookieStore) *ViewController {
	return &ViewController{
		Session: session,
		Service: service,
	}
}

func (vc *ViewController) NotFoundViewHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		store, _ := vc.Session.Get(r, "default")
		user := store.Values["user"]
		if user != nil {
			view.NotFoundPage(user.(*dto.UserDTO)).Render(r.Context(), w)
			return
		}

		view.NotFoundPage(nil).Render(r.Context(), w)
	}
}

func (vc *ViewController) SignInHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		store, _ := vc.Session.Get(r, "default")
		dto := &dto.SignInPageDTO{}

		for _, flash := range store.Flashes("error") {
			dto.Error = flash.(string)
		}

		sessions.Save(r, w)
		view.SignInPage(dto).Render(r.Context(), w)
	}
}

func (vc *ViewController) SignInPostHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := vc.Session.Get(r, "default")
		decoder := schema.NewDecoder()

		var payload dto.UserSignInRequestDTO

		err := r.ParseForm()
		if err != nil {
			log.Println("error parsing form:", err.Error())
			http.Redirect(w, r, "/auth/sign-in", http.StatusMovedPermanently)
			return
		}

		err = decoder.Decode(&payload, r.PostForm)
		if err != nil {
			log.Println("error decoding payload: ", err.Error())
			http.Redirect(w, r, "/auth/sign-in", http.StatusMovedPermanently)
			return
		}

		if strings.TrimSpace(payload.Username) == "" || strings.TrimSpace(payload.Password) == "" {
			session.AddFlash("Username or password invalid, please try again.", "error")
			sessions.Save(r, w)
			http.Redirect(w, r, "/auth/sign-in", http.StatusMovedPermanently)
			return
		}

		user, err := vc.Service.AuthenticationService.SignIn(r.Context(), payload.Username, payload.Password)
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
			http.Redirect(w, r, "/auth/sign-in", http.StatusMovedPermanently)
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

func (vc *ViewController) RegisterHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		store, _ := vc.Session.Get(r, "default")

		var payload dto.RegisterUserRequestDTO
		decoder := schema.NewDecoder()

		r.ParseForm()
		err := decoder.Decode(&payload, r.PostForm)
		if err == nil {
			payload := &dto.RegisterPageDTO{
				RegisterFieldDTO: &dto.RegisterFieldDTO{
					Username: struct {
						Value  string
						Errors []string
					}{
						Value: payload.Username,
					},
					Email: struct {
						Value  string
						Errors []string
					}{Value: payload.Email},
					PasswordErrors:        []string{},
					ConfirmPasswordErrors: []string{},
				},
				Error: "",
			}

			// redirect result check for flashes

			for _, flash := range store.Flashes("username") {
				payload.RegisterFieldDTO.Username.Errors = append(payload.RegisterFieldDTO.Username.Errors, flash.(string))
			}

			for _, flash := range store.Flashes("email") {
				payload.RegisterFieldDTO.Email.Errors = append(payload.RegisterFieldDTO.Email.Errors, flash.(string))
			}

			for _, flash := range store.Flashes("password") {
				payload.RegisterFieldDTO.PasswordErrors = append(payload.RegisterFieldDTO.PasswordErrors, flash.(string))
			}

			for _, flash := range store.Flashes("confirm_password") {
				payload.RegisterFieldDTO.ConfirmPasswordErrors = append(payload.RegisterFieldDTO.ConfirmPasswordErrors, flash.(string))
			}

			for _, flash := range store.Flashes("error") {
				payload.Error = flash.(string)
			}

			err := sessions.Save(r, w)
			if err != nil {
				log.Println("error saving session :", err.Error())
			}
			view.RegisterPage(payload).Render(r.Context(), w)
			return
		}

		view.RegisterPage(nil).Render(r.Context(), w)
	}
}

func (vc *ViewController) HomepageViewHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		homeDTO := &dto.HomepageDTO{
			Posts: []dto.PostDTO{},
		}

		store, _ := vc.Session.Get(r, "default")
		userStore := store.Values["user"]
		if userStore != nil {
			homeDTO.User = userStore.(*dto.UserDTO)
		}

		view.Homepage(homeDTO).Render(r.Context(), w)
	}
}

func (vc *ViewController) CreateNewPostViewHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		store, _ := vc.Session.Get(r, "default")
		user := store.Values["user"].(*dto.UserDTO)

		createNewPostDTO := &dto.CreateNewPostDTO{
			User: user,
		}

		view.CreateNewPostPage(createNewPostDTO).Render(r.Context(), w)
	}
}

func (vc *ViewController) CreatePostHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (vc *ViewController) ProfileHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		store, _ := vc.Session.Get(r, "default")
		user := store.Values["user"].(*dto.UserDTO)

		profileDTO := &dto.ProfilePageDTO{
			User:  user,
			Posts: []dto.PostDTO{},
		}

		for i := 0; i < 25; i++ {
			profileDTO.Posts = append(profileDTO.Posts, dto.PostDTO{
				ID:          strconv.Itoa(i + 1),
				Content:     fmt.Sprintf("this is my %d post", i+1),
				Replies:     []dto.PostDTO{},
				ReplyCounts: i,
				Likes:       i,
				SavedCounts: i,
				Impressions: i,
				Poster: dto.UserDTO{
					Username: user.Username,
				},
				PostedAt: time.Now(),
			})
		}

		view.ProfilePage(profileDTO).Render(r.Context(), w)
	}
}

func (vc *ViewController) SettingsHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		store, _ := vc.Session.Get(r, "default")
		user := store.Values["user"].(*dto.UserDTO)

		view.SettingsPage(&dto.SettingsPageDto{
			User: user,
		}).Render(r.Context(), w)
	}
}
