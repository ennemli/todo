package routing

import (
	"github.com/ennemli/todo/user/configs"
	handlers "github.com/ennemli/todo/user/internal/handlers/user"
	"github.com/ennemli/todo/user/internal/models/user"
	"github.com/ennemli/todo/user/internal/server"
	"github.com/go-chi/chi/v5"
)

func InitRouting(store user.Store) {
	r := server.GetRouter()
	userHandler := handlers.NewUserHandler(store)
	r.Route(configs.GetConfig().Service.API_ENDPOINT, func(r chi.Router) {
		r.Get("/", userHandler.GetUsers)
		r.Post("/", userHandler.CreateUser)
		r.Route("/{id:^[0-9]+$}", func(r chi.Router) {
			r.Get("/", userHandler.GetUserById)
			r.Delete("/", userHandler.DeleteUserById)
			r.Put("/", userHandler.UpdateUser)
		})
		r.Get("/{name:^[a-zA-Z][a-zA-Z-]+$}", userHandler.GetUserByName)
	})
}
