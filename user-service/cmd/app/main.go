package main

import (
	"time"

	"github.com/ennemli/todo/user/internal/db"
	"github.com/ennemli/todo/user/internal/middlewares"
	"github.com/ennemli/todo/user/internal/models/user"
	"github.com/ennemli/todo/user/internal/routing"
	"github.com/ennemli/todo/user/internal/server"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func main() {
	server.InitServerEnv()
	s := server.NewServer()
	r := server.GetRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.SetTimeOut(time.Second * 2))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	db.GetDB().AutoMigrate(&user.User{})
	routing.InitRouting(user.NewStore())
	s.ListenAndServe()
}
