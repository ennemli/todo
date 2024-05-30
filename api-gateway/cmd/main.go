package main

import (
	"github.com/ennemli/apigateway/internal/middlewares"
	"github.com/ennemli/apigateway/internal/proxy"
	"github.com/ennemli/apigateway/internal/server"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	server.InitServerEnv()
	s := server.NewServer()
	r := server.GetRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.With(middlewares.WithAuth).Mount("/api/todos", proxy.TodoAPIProxy())
	r.With(middlewares.WithAuth).Mount("/api/users", proxy.UsersAPIProxy())
	r.Mount("/auth", proxy.AuthAPIProxy())
	s.ListenAndServe()
}
