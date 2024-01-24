package service

type Service struct {
	UserService *UserService
}

func NewService() *Service {
	return &Service{
		UserService: NewUserService(),
	}
}
