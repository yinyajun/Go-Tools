package admin

import (
	"fmt"
	"net/http"

	"github.com/yinyajun/Golang-Toys/pepper_cache/cache/server/core"
	"github.com/yinyajun/Golang-Toys/pepper_cache/cache/server/setting"
)

type Server struct {
	*core.Server
	AdminPort int
}

func (s *Server) HttpListen() {
	http.Handle("/cluster", s.clusterHandler())
	http.Handle("/stat", s.statusHandler())
	http.ListenAndServe(fmt.Sprintf(":%d", s.AdminPort), nil)
}

func NewAdminServer(s *core.Server) *Server {
	return &Server{s, setting.AdminPort}
}
