package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/ennemli/todo/user/configs"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	httpServer *http.Server
	router     *chi.Mux
}

var server *Server = nil

func NewServer() *Server {
	if server == nil {
		r := chi.NewRouter()
		server = &Server{
			router: r,
			httpServer: &http.Server{
				Addr:    fmt.Sprintf(":%d", configs.GetConfig().Service.APP_PORT),
				Handler: r,
			},
		}
	}
	return server
}

func (s *Server) ListenAndServe() {
	l, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		panic(err)
	}
	s.httpServer.Serve(l)
}

func InitServerEnv() {
	configs.Initialize(".env", ".", "env")
}

func GetRouter() *chi.Mux {
	return server.router
}

func GetServer() *http.Server {
	return server.httpServer
}
