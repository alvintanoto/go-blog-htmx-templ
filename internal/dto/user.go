package dto

type UserDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type RegisterUserRequestDTO struct {
	Username        string `schema:"username"`
	Email           string `schema:"email"`
	Password        string `schema:"password"`
	ConfirmPassword string `schema:"confirm_password"`
}

type UserSignInRequestDTO struct {
	Username string `schema:"username"`
	Password string `schema:"password"`
}
