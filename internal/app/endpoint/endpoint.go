package endpoint

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Service interface {
	PerformPull(repoPath string, privateKeyPath string) error
	RestartContainers(workingDir string, composeFilePath string) error
}

type Endpoint struct {
	service Service
}

func New(service Service) *Endpoint {
	return &Endpoint{
		service: service,
	}
}

func (e *Endpoint) Status(ctx echo.Context) error {
	// TODO: perform health checks
	err := ctx.JSON(http.StatusOK, "OK")
	if err != nil {
		return err
	}
	return nil
}

func (e *Endpoint) Update(ctx echo.Context) error {
	err := e.service.PerformPull("/repo", "/.ssh/id_rsa")
	if err != nil {
		return err
	}
	err = e.service.RestartContainers("/repo", "/repo/docker-compose.yml")
	if err != nil {
		return err
	}
	err = ctx.JSON(http.StatusOK, "OK")
	// TODO: tag the commit as deployed
	return nil
}
