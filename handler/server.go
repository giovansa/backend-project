package handler

import (
	"github.com/SawitProRecruitment/UserService/internal"
	"github.com/SawitProRecruitment/UserService/repository"
)

type Server struct {
	Cfg        internal.Config
	Repository repository.RepositoryInterface
}

type NewServerOptions struct {
	Repository repository.RepositoryInterface
}

func NewServer(cfg internal.Config, opts NewServerOptions) *Server {
	return &Server{
		Cfg:        cfg,
		Repository: opts.Repository,
	}
}
