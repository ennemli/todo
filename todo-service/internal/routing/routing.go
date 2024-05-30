package routing

import (
	"github.com/ennemli/todo/todo/configs"
	handlers "github.com/ennemli/todo/todo/internal/handlers/todo"
	"github.com/ennemli/todo/todo/internal/models/todo"
	"github.com/ennemli/todo/todo/internal/server"
	"github.com/go-chi/chi/v5"
)

func InitRouting(store todo.Store) {
	r := server.GetRouter()
	todoHandlers := handlers.NewTodoHandlers(store)
	r.Route(configs.GetConfig().Service.API_ENDPOINT, func(r chi.Router) {
		r.Get("/", todoHandlers.GetTodos)
		r.Post("/", todoHandlers.CreateTodo)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", todoHandlers.GetTodoById)
			r.Delete("/", todoHandlers.DeleteTodoById)
			r.Put("/", todoHandlers.UpdateTodo)
		})
	})
}
