package dto

type HomepageDTO struct {
	User  *UserDTO  `json:"user"`
	Posts []PostDTO `json:"posts"`
	Error string
}

type CreateNewPostDTO struct {
	User    *UserDTO `json:"user"`
	Content string   `json:"content"`
	Error   string
}

type RegisterFieldDTO struct {
	Username struct {
		Value  string
		Errors []string
	}

	Email struct {
		Value  string
		Errors []string
	}

	PasswordErrors        []string
	ConfirmPasswordErrors []string
}

type SignInPageDTO struct {
	Error string
}

type RegisterPageDTO struct {
	RegisterFieldDTO *RegisterFieldDTO
	Error            string
}

type DraftPageDTO struct {
	User  *UserDTO  `json:"user"`
	Posts []PostDTO `json:"posts"`
	Error string
}

type SettingsPageDto struct {
	User  *UserDTO `json:"user"`
	Error string
}

type ProfilePageDTO struct {
	User  *UserDTO  `json:"user"`
	Posts []PostDTO `json:"posts"`
	Error string
}

type PostDetailDTO struct {
	User  *UserDTO  `json:"user"`
	Posts []PostDTO `json:"posts"`
	Error string
}
