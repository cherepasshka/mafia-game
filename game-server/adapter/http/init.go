package http_mafia

import (
	// "net/http"

	"github.com/go-chi/chi"

	usersdb "soa.mafia-game/game-server/domain/models/storage"
)

type HttpHandler struct {
	users *usersdb.Storage
}

func New(users *usersdb.Storage) *HttpHandler {
	return &HttpHandler{
		users: users,
	}
}

func (handler *HttpHandler) AddHandlers(router *chi.Mux) {
	router.Route("/users/{login}", func(r chi.Router) {
		r.Post("/", handler.CreateUser)
		r.Get("/", handler.GetUser)
		r.Put("/", handler.UpdateUser)
		r.Delete("/", handler.DeleteUser)
	})
	router.HandleFunc("/users/", handler.GetUsers)
}
