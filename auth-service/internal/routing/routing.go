package routing

import (
	"github.com/ennemli/todo/todo/configs"
	handlers "github.com/ennemli/todo/todo/internal/handlers/auth"
	"github.com/ennemli/todo/todo/internal/server"
	"github.com/go-chi/chi/v5"
)

func InitRouting() {
	r := server.GetRouter()

	authHandler := handlers.NewAuthHandler(new(handlers.UserServiceClient))
	r.Route(configs.GetConfig().Service.API_ENDPOINT, func(r chi.Router) {
		r.Post("/", authHandler.LoginHandler)
		r.Post("/valid", authHandler.ValidateHandler)
	})
}
