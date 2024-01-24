package dto

type HomepageDTO struct {
	User  *UserDTO  `json:"user"`
	Posts []PostDTO `json:"posts"`
}

type CreateNewPostDTO struct {
	User *UserDTO `json:"user"`
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

type RegisterPageDTO struct {
	RegisterFieldDTO *RegisterFieldDTO
}
