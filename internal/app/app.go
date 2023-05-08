package app

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/vi350/webhook-deployer/internal/app/endpoint"
	"github.com/vi350/webhook-deployer/internal/app/mw"
	"github.com/vi350/webhook-deployer/internal/app/service"
)

type App struct {
	e    *endpoint.Endpoint
	s    *service.Service
	echo *echo.Echo
}

func New() (*App, error) {
	a := &App{}
	a.s = service.New()
	a.e = endpoint.New(a.s)

	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	a.echo = echo.New()
	a.echo.GET("/status", a.e.Status)
	a.echo.POST("/update", a.e.Update, mw.AuthorizePushEvent)
	return a, nil
}

func (a *App) Run() error {
	err := a.echo.Start(":80")
	if err != nil {
		return err
	}
	return nil
}
