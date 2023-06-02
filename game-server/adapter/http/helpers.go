package http_mafia

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

// GET /users/[id]
func (handler *HttpHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	log.Printf("got %v", id)
}

// GET /users/[id1]....
func (handler *HttpHandler) GetUsers(w http.ResponseWriter, r *http.Request) {

}
