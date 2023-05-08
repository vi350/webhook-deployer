package endpoint

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
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
	fmt.Println("update started")
	err := e.service.PerformPull("/repo", "/root/.ssh/id_rsa")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("pulled")
	err = e.service.RestartContainers("/repo", "/repo/docker-compose.yml")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("restarted")
	err = ctx.JSON(http.StatusOK, "OK")
	// TODO: tag the commit as deployed
	return nil
}
