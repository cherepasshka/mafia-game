package http_mafia

import (
	// "net/http"

	"github.com/go-chi/chi"

	usersdb "soa.mafia-game/game-server/domain/models/users_db"
)

type HttpHandler struct {
	users *usersdb.UsersStorage
}

func New(users *usersdb.UsersStorage) *HttpHandler {
	return &HttpHandler{
		users: users,
	}
}

func (handler *HttpHandler) AddHandlers(router *chi.Mux) {
	router.Route("/users/{login}", func(r chi.Router) {
		// r.Post("/", handler.GetUser)
		r.Get("/", handler.GetUser)
		// r.Put("/", UpdateUser)
		// r.Delete("/", DeleteUser)
	})
}
