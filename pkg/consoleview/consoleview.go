package consoleview

import (
	"log"

	"example.com/pkg/domain"
)

type CMDApp struct {
	service  domain.UserService
	handlers map[string]Handler
}

type Handler struct {
	Name       string
	HandleFunc func() error
}

func NewCMDApp(service domain.UserService) (*CMDApp, error) {
	app := CMDApp{service: service}
	app.handlers = make(map[string]Handler, 5)
	app.initRoutes()
	return &app, nil
}

func (app *CMDApp) Route(name string, f func() error) {
	app.handlers[name] = Handler{Name: name, HandleFunc: f}
}

func (app *CMDApp) Serve(h Handler) {
	err := h.HandleFunc()
	if err != nil {
		log.Fatalln(err)
	}
}

func (app *CMDApp) Run() error {
	return app.handlers["auth"].HandleFunc()
}
