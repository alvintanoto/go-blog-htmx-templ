package controller

type Controller struct {
	ViewController *ViewController
}

func NewController() *Controller {
	return &Controller{
		ViewController: NewViewController(),
	}
}
