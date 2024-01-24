package dto

type UserDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type RegisterUserRequestDTO struct {
	Username        string `schema:"username"`
	Email           string `schema:"email"`
	Password        string `schema:"password"`
	ConfirmPassword string `schema:"confirm_password"`
}
