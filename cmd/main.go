package main

import (
	"fmt"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/internal"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/transport"
	"log"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	var server handler.HandlerInterface = newServer()

	transport.RegisterHandler(e, server)
	e.Logger.Fatal(e.Start(":8080"))
}

func newServer() *handler.Server {
	cfg, err := internal.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
			"postgres",
			cfg.DB.User,
			cfg.DB.Password,
			cfg.DB.Host,
			cfg.DB.Port,
			cfg.DB.DatabaseName),
	})
	opts := handler.NewServerOptions{
		Repository: repo,
	}
	return handler.NewServer(cfg, opts)
}
