package app

import (
	"github.com/EnduranNSU/end-user-info/internal/domain"
)

type Server struct {
	UserInfoRepository *domain.UserInfoRepository
}

func (s *Server) StartServer() {

}

func SetupServer(UserInfoRepository *domain.UserInfoRepository) *Server {
	return &Server{
		UserInfoRepository: UserInfoRepository,
	}
}
