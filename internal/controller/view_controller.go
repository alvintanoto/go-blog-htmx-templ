package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"alvintanoto.id/blog-htmx-templ/internal/dto"
	"alvintanoto.id/blog-htmx-templ/internal/repository"
	"alvintanoto.id/blog-htmx-templ/internal/service"
	vcomponent "alvintanoto.id/blog-htmx-templ/internal/view/component"
	verror "alvintanoto.id/blog-htmx-templ/internal/view/error"
	vpages "alvintanoto.id/blog-htmx-templ/internal/view/page"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/rbcervilla/redisstore/v9"
	"golang.org/x/crypto/bcrypt"
)

type ViewController struct {
	Session *redisstore.RedisStore
	Service *service.Service
}

func NewViewController(service *service.Service, session *redisstore.RedisStore) *ViewController {
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
			verror.NotFound(user.(*dto.UserDTO)).Render(r.Context(), w)
			return
		}

		verror.NotFound(nil).Render(r.Context(), w)
	}
}

func (vc *ViewController) SignInHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		store, _ := vc.Session.Get(r, "default")
		data := &dto.SignInPageDTO{}

		switch r.Method {
		case http.MethodGet:
			vpages.SignIn(data).Render(r.Context(), w)
		case http.MethodPost:
			if flashes := store.Flashes("error"); len(flashes) > 0 {
				data.Error = flashes[len(flashes)-1].(string)

				store.Save(r, w)
				vpages.SignIn(data).Render(r.Context(), w)
				return
			}

			decoder := schema.NewDecoder()

			var payload dto.UserSignInRequestDTO

			err := r.ParseForm()
			if err != nil {
				log.Println("error parsing form:", err.Error())
				http.Redirect(w, r, "/auth/sign-in", http.StatusTemporaryRedirect)
				return
			}

			err = decoder.Decode(&payload, r.PostForm)
			if err != nil {
				log.Println("error decoding payload: ", err.Error())
				http.Redirect(w, r, "/auth/sign-in", http.StatusTemporaryRedirect)
				return
			}

			if strings.TrimSpace(payload.Username) == "" || strings.TrimSpace(payload.Password) == "" {
				store.AddFlash("Username or password invalid, please try again.", "error")
				store.Save(r, w)
				http.Redirect(w, r, "/auth/sign-in", http.StatusTemporaryRedirect)
				return
			}

			user, err := vc.Service.AuthenticationService.SignIn(r.Context(), payload.Username, payload.Password)
			if err != nil {
				switch err {
				case repository.ErrorRecordNotFound,
					bcrypt.ErrMismatchedHashAndPassword:
					store.AddFlash("Username or password invalid, please try again.", "error")
				default:
					store.AddFlash("Failed to sign in, please try again later.", "error")
				}
				// TODO: redirect to sign in with flash error
				store.Save(r, w)
				http.Redirect(w, r, "/auth/sign-in", http.StatusTemporaryRedirect)
				return
			}

			store.Values["user"] = &dto.UserDTO{
				ID:        user.ID,
				Username:  user.Username,
				Email:     user.Email,
				CreatedAt: user.CreatedAt,
			}
			store.Save(r, w)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}

func (vc *ViewController) RegisterHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := vc.Session.Get(r, "default")
		registerPageDataDTO := &dto.RegisterPageDTO{
			RegisterFieldDTO: &dto.RegisterFieldDTO{
				Username: struct {
					Value  string
					Errors []string
				}{
					Value: "",
				},
				Email: struct {
					Value  string
					Errors []string
				}{Value: ""},
				PasswordErrors:        []string{},
				ConfirmPasswordErrors: []string{},
			},
			Error: "",
		}

		switch r.Method {
		case http.MethodGet:
			vpages.Register(registerPageDataDTO).Render(r.Context(), w)
			return
		case http.MethodPost:
			var payload dto.RegisterUserRequestDTO
			decoder := schema.NewDecoder()

			err := r.ParseForm()
			if err != nil {
				log.Println("error parsing form: ", err.Error())
				http.Redirect(w, r, "/auth/register", http.StatusTemporaryRedirect)
				return
			}
			err = decoder.Decode(&payload, r.PostForm)
			if err != nil {
				log.Println("error decoding payload: ", err.Error())
				http.Redirect(w, r, "/auth/register", http.StatusTemporaryRedirect)
				return
			}

			if flashes := session.Flashes("validation_error"); len(flashes) > 0 {
				registerPageDataDTO.RegisterFieldDTO.Username.Value = payload.Username
				registerPageDataDTO.RegisterFieldDTO.Email.Value = payload.Email

				// redirect result check for flashes
				for _, flash := range session.Flashes("username") {
					registerPageDataDTO.RegisterFieldDTO.Username.Errors = append(registerPageDataDTO.RegisterFieldDTO.Username.Errors, flash.(string))
				}

				for _, flash := range session.Flashes("email") {
					registerPageDataDTO.RegisterFieldDTO.Email.Errors = append(registerPageDataDTO.RegisterFieldDTO.Email.Errors, flash.(string))
				}

				for _, flash := range session.Flashes("password") {
					registerPageDataDTO.RegisterFieldDTO.PasswordErrors = append(registerPageDataDTO.RegisterFieldDTO.PasswordErrors, flash.(string))
				}

				for _, flash := range session.Flashes("confirm_password") {
					registerPageDataDTO.RegisterFieldDTO.ConfirmPasswordErrors = append(registerPageDataDTO.RegisterFieldDTO.ConfirmPasswordErrors, flash.(string))
				}

				for _, flash := range session.Flashes("error") {
					registerPageDataDTO.Error = flash.(string)
				}

				sessions.Save(r, w)
				vpages.Register(registerPageDataDTO).Render(r.Context(), w)
				return
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

			if len(username) <= 6 || len(username) > 25 {
				session.AddFlash("Username length must be more than 6 character and less than 25 character.", "username")
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
				session.AddFlash("1", "validation_error")
				err := sessions.Save(r, w)
				if err != nil {
					log.Println("err saving session:", err.Error())
					http.Redirect(w, r, "/auth/register", http.StatusSeeOther)
					return
				}
				http.Redirect(w, r, "/auth/register", http.StatusTemporaryRedirect)
				return
			}

			user, err := vc.Service.UserService.RegisterUser(username, email, password)
			if err != nil {
				log.Println("failed registering new user: ", err.Error())
				switch err {
				case repository.ErrorConstraintViolation:
					session.AddFlash("Username already used, please try another username.", "error")
				default:
					session.AddFlash("Failed registering new user, please try again later.", "error")
				}
				sessions.Save(r, w)
				http.Redirect(w, r, "/register", http.StatusTemporaryRedirect)
				return
			}

			session.Values["user"] = &dto.UserDTO{
				ID:        user.ID,
				Username:  user.Username,
				Email:     user.Email,
				CreatedAt: user.CreatedAt,
			}
			sessions.Save(r, w)

			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}

func (vc *ViewController) HomepageViewHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		store, _ := vc.Session.Get(r, "default")

		homeDTO := &dto.HomepageDTO{
			Posts: []dto.PostDTO{},
		}

		userStore := store.Values["user"]

		if userStore != nil {
			homeDTO.User = userStore.(*dto.UserDTO)

			// TODO: get user timeline posts
			posts, err := vc.Service.PostService.GetHomeTimeline(homeDTO.User, 0)
			if err != nil {
				homeDTO.Error = "Failed to get timeline, please try again later"
			}

			homeDTO.Posts = posts
		}

		vpages.Home(homeDTO).Render(r.Context(), w)
	}
}

func (vc *ViewController) HomepageInfiniteScrollHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		store, _ := vc.Session.Get(r, "default")
		user := store.Values["user"].(*dto.UserDTO)

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))

		// TODO: get user timeline posts
		posts, err := vc.Service.PostService.GetUserPosts(user, page)
		if err != nil {
			return
		}

		nPage := page + 1
		vcomponent.Posts(posts, fmt.Sprintf("/profile/load-more-posts?page=%d", nPage)).Render(r.Context(), w)
	}
}

func (vc *ViewController) CreatePostHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := vc.Session.Get(r, "default")
		user := session.Values["user"].(*dto.UserDTO)

		id := r.URL.Query().Get("id")

		createPostDTO := &dto.CreateNewPostDTO{
			User: user,
		}

		if flashes := session.Flashes("error"); len(flashes) > 0 {
			createPostDTO.Error = flashes[len(flashes)-1].(string)

			session.Save(r, w)
			vpages.CreatePost(createPostDTO).Render(r.Context(), w)
			return
		}

		switch r.Method {
		case http.MethodGet:
			if id != "" {
				post, err := vc.Service.PostService.GetUserPost(user, id)
				if err != nil {
					// failed just redirect to drafts
					http.Redirect(w, r, "/draft/", http.StatusSeeOther)
					return
				}

				createPostDTO.Content = post.Content
			}
			vpages.CreatePost(createPostDTO).Render(r.Context(), w)
		case http.MethodPost:
			decoder := schema.NewDecoder()

			var payload dto.SubmitPostDTO

			err := r.ParseForm()
			if err != nil {
				log.Println("error parsing form: ", err.Error())
				http.Redirect(w, r, "/post/new-post", http.StatusMovedPermanently)
				return
			}
			err = decoder.Decode(&payload, r.PostForm)
			if err != nil {
				log.Println("error decoding payload: ", err.Error())
				http.Redirect(w, r, "/post/new-post", http.StatusMovedPermanently)
				return
			}

			content := strings.TrimSpace(payload.Content)
			if len(content) <= 0 {
				log.Println("create post handler error: empty content value")

				session.AddFlash("Content cannot be empty.", "error")
				sessions.Save(r, w)

				http.Redirect(w, r, "/post/new-post", http.StatusMovedPermanently)
				return
			}

			// data not exist in the database
			err = vc.Service.PostService.CreatePost(user.ID, id, content, payload.SubmitType)
			if err != nil {
				log.Println("error posting / drafting post: ", err.Error())

				session.AddFlash("failed to create / draft a new post, please try again later.", "error")
				sessions.Save(r, w)

				http.Redirect(w, r, "/post/new-post", http.StatusMovedPermanently)
				return
			}

			if payload.SubmitType == "draft" {
				http.Redirect(w, r, "/draft/", http.StatusSeeOther)
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}

func (vc *ViewController) PostDetailHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		store, _ := vc.Session.Get(r, "default")
		user := store.Values["user"].(*dto.UserDTO)

		pathParam := mux.Vars(r)
		postID := pathParam["id"]
		if postID == "" {
			verror.NotFound(user).Render(r.Context(), w)
			return
		}

		fmt.Println(postID)
		// post, err := vc.Service.PostService.GetPostDetail(user, postID)
		// if err != nil {

		// }

		vpages.PostDetail(dto.PostDetailDTO{
			User: user,
		}).Render(r.Context(), w)
		return
	}
}

func (vc *ViewController) DraftHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		store, _ := vc.Session.Get(r, "default")
		user := store.Values["user"].(*dto.UserDTO)

		draftDTO := &dto.DraftPageDTO{
			User:  user,
			Posts: []dto.PostDTO{},
		}

		posts, err := vc.Service.PostService.GetUserDraft(user, 0)
		if err != nil {
			draftDTO.Error = "Failed to get user drafts, please try again later"
			vpages.Draft(draftDTO).Render(r.Context(), w)
			return
		}

		draftDTO.Posts = posts
		vpages.Draft(draftDTO).Render(r.Context(), w)
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

		posts, err := vc.Service.PostService.GetUserPosts(user, 0)
		if err != nil {
			profileDTO.Error = "Failed to get user profile post, please try again later"
			vpages.Profile(profileDTO).Render(r.Context(), w)
			return
		}

		profileDTO.Posts = posts
		vpages.Profile(profileDTO).Render(r.Context(), w)
	}
}

func (vc *ViewController) ProfilePostInfiniteScrollHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		store, _ := vc.Session.Get(r, "default")
		user := store.Values["user"].(*dto.UserDTO)

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))

		posts, err := vc.Service.PostService.GetUserPosts(user, page)
		if err != nil {
			return
		}

		nPage := page + 1
		vcomponent.Posts(posts, fmt.Sprintf("/profile/load-more-posts?page=%d", nPage)).Render(r.Context(), w)
	}
}

func (vc *ViewController) SettingsHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		store, _ := vc.Session.Get(r, "default")
		user := store.Values["user"].(*dto.UserDTO)

		vpages.Settings(&dto.SettingsPageDto{
			User: user,
		}).Render(r.Context(), w)
	}
}

func (vc *ViewController) SignOutHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			store, _ := vc.Session.Get(r, "default")
			store.Values["user"] = nil
			store.Save(r, w)

			w.Header().Set("Hx-redirect", "/")
			w.Write([]byte(""))
		}
	}
}

func (vc *ViewController) HideModal() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(""))
	}
}

func (vc *ViewController) ShowSignOutConfirmation() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vcomponent.SignOutModal().Render(r.Context(), w)
	}
}
