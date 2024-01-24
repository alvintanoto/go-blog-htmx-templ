package dto

type HomepageDTO struct {
	User  *UserDTO  `json:"user"`
	Posts []PostDTO `json:"posts"`
}

type CreateNewPostDTO struct {
	User *UserDTO `json:"user"`
}
