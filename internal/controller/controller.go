package controller

type Controller struct {
	Middlewares    *Middlewares
	ViewController *ViewController
}

func NewController() *Controller {
	return &Controller{
		Middlewares:    NewMiddleware(),
		ViewController: NewViewController(),
	}
}
