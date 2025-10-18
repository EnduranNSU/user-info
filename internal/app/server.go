package app

import (
	"net/http"
	"time"

	httpin "github.com/EnduranNSU/end-user-info/internal/adapter/in/http"
	svcuserinfo "github.com/EnduranNSU/end-user-info/internal/service"
)

type Server struct {
	Svc  svcuserinfo.Service
	Addr string
}

func SetupServer(svc svcuserinfo.Service, addr string) *Server {
	return &Server{Svc: svc, Addr: addr}
}

func (s *Server) StartServer() error {
	h := httpin.NewUserInfoHandler(s.Svc)
	engine := httpin.NewGinRouter(h)

	srv := &http.Server{
		Addr:              s.Addr,
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}
	return srv.ListenAndServe()
}
