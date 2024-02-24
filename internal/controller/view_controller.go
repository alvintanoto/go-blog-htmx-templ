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

			username := r.PostFormValue("username")
			password := r.PostFormValue("password")

			if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
				store.AddFlash("Username or password invalid, please try again.", "error")
				store.Save(r, w)
				http.Redirect(w, r, "/auth/sign-in", http.StatusTemporaryRedirect)
				return
			}

			user, err := vc.Service.AuthenticationService.SignIn(r.Context(), username, password)
			if err != nil {
				switch err {
				case repository.ErrorRecordNotFound,
					bcrypt.ErrMismatchedHashAndPassword:
					store.AddFlash("Username or password invalid, please try again.", "error")
				default:
					store.AddFlash("Failed to sign in, please try again later.", "error")
				}

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

		if flashes := session.Flashes("error"); len(flashes) > 0 {
			registerPageDataDTO.Error = flashes[len(flashes)-1].(string)

			session.Save(r, w)
			vpages.Register(registerPageDataDTO).Render(r.Context(), w)
			return
		}

		switch r.Method {
		case http.MethodGet:
			vpages.Register(registerPageDataDTO).Render(r.Context(), w)
			return
		case http.MethodPost:
			username := strings.TrimSpace(r.PostFormValue("username"))
			email := strings.TrimSpace(r.PostFormValue("email"))
			password := r.PostFormValue("password")
			confirmPassword := r.PostFormValue("confirm_password")
			errCount := 0

			if flashes := session.Flashes("validation_error"); len(flashes) > 0 {
				registerPageDataDTO.RegisterFieldDTO.Username.Value = username
				registerPageDataDTO.RegisterFieldDTO.Email.Value = password

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
				http.Redirect(w, r, "/auth/register", http.StatusTemporaryRedirect)
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

		}

		vpages.Home(homeDTO).Render(r.Context(), w)
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
		post, err := vc.Service.PostService.GetPostDetail(postID)
		if err != nil {
			http.Redirect(w, r, fmt.Sprintf("/"), http.StatusSeeOther)
			return
		}

		vpages.PostDetail(dto.PostDetailDTO{
			User: user,
			Post: *post,
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

		profileDTO := &dto.PageDTO{
			RouteName: "Profile",
			User:      user,
		}

		vpages.Profile(profileDTO).Render(r.Context(), w)
	}
}

func (vc *ViewController) ProfilePostInfiniteScrollHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		store, _ := vc.Session.Get(r, "default")
		user := store.Values["user"].(*dto.UserDTO)

		lastPosition, _ := strconv.Atoi(r.URL.Query().Get("last_position"))

		posts, err := vc.Service.PostService.GetUserPosts(user, lastPosition)
		if err != nil {
			return
		}

		newLastPositionID := 0
		if len(posts) > 0 {
			newLastPositionID = posts[len(posts)-1].ID
		}

		vcomponent.Posts(posts, newLastPositionID).Render(r.Context(), w)
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

func (vc *ViewController) Populate() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		store, _ := vc.Session.Get(r, "default")
		user := store.Values["user"].(*dto.UserDTO)

		go func() {
			var contents []string

			for i := 0; i < 1000000; i++ {
				contents = append(contents, fmt.Sprintf("[%d] %s", i+1, `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam viverra ut lectus vel tincidunt. In efficitur nisi ultricies est tempus, non aliquet diam vulputate. Sed ullamcorper, nulla eget ullamcorper elementum, ligula nulla ornare augue, eget convallis orci lacus ac est. Etiam justo nulla, tincidunt et nisl ac, volutpat vestibulum nunc. Quisque ac lacus eu tortor mattis porta ac sit amet quam. Pellentesque ultrices pulvinar aliquam. Vestibulum eget quam leo. Sed sed lectus vitae metus placerat fringilla.
			Aenean vitae justo vitae lectus auctor aliquet ut et ante. Mauris tempus vehicula nisi nec varius. Nam enim nunc, suscipit sit amet tristique ut, tincidunt eu leo. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Phasellus pellentesque efficitur sapien, sed hendrerit ligula egestas eu. Phasellus turpis dui, imperdiet eget neque at, bibendum vulputate arcu. Duis at pretium felis. Donec in urna eget felis lobortis dapibus. Praesent tempor lorem libero, id vestibulum ipsum suscipit sed. Duis varius urna et elit venenatis placerat. Mauris vitae enim id ante semper blandit at nec justo.
			In tristique enim id odio rutrum, at ultrices tortor consequat. Aenean congue tincidunt interdum. Integer quis urna lacinia, mollis ipsum id, lobortis augue. Proin quis suscipit nibh. Morbi sit amet iaculis ante. Mauris sit amet lacinia nunc, ut commodo est. Quisque a vulputate quam, nec placerat elit.
			Aliquam quis nulla eget sem pretium congue vitae vitae mi. Curabitur sagittis ex ut ex gravida efficitur. Ut placerat vulputate metus in volutpat. Nunc nec lobortis orci, sit amet vehicula tellus. Morbi non ligula at est sollicitudin dapibus. Aenean quis magna justo. Aenean sed orci vel tellus interdum efficitur. Sed mi libero, scelerisque vel lobortis in, finibus dapibus mauris. Maecenas sit amet aliquet ligula. Etiam in urna at nunc rhoncus pulvinar. Aliquam pretium molestie metus.
			Duis pharetra metus eu tristique eleifend. Proin hendrerit interdum mauris non gravida. Etiam malesuada purus dui, sit amet elementum ipsum fringilla nec. Vivamus mauris urna, faucibus a ullamcorper sed, tincidunt nec nisi. Vestibulum at gravida augue. Sed sed condimentum erat. Donec finibus placerat augue, sit amet varius risus mollis a. Etiam orci erat, posuere sit amet volutpat non, volutpat vitae lacus. In ut pretium massa, at consequat magna. Pellentesque sit amet posuere risus. Phasellus interdum nulla vitae lectus blandit scelerisque. Vestibulum sit amet turpis vulputate, varius ante tempus, rutrum arcu. Nulla facilisi.
			`))
			}

			vc.Service.PostService.Populate(user.ID, contents)
		}()

		w.Write([]byte("success"))
	}
}
