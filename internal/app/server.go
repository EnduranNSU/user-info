package app

import (
	"net/http"
	"time"

	httpin "github.com/EnduranNSU/end-user-info/internal/adapter/in/http"
	"github.com/EnduranNSU/end-user-info/internal/domain"
)

type Server struct {
	Repo domain.UserInfoRepository
	Addr string
}

func SetupServer(repo domain.UserInfoRepository) *Server {
	return &Server{Repo: repo, Addr: ":8080"}
}

func (s *Server) StartServer() error {
	h := httpin.NewUserInfoHandler(s.Repo)
	engine := httpin.NewGinRouter(h)

	srv := &http.Server{
		Addr:              s.Addr,
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}
	return srv.ListenAndServe()
}
