package http_mafia

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

// GET /users/{login}
func (handler *HttpHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	login := chi.URLParam(r, "login")
	log.Printf("got %v", login)
	user, exists := handler.users.Get(login)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
	}
	jsonValue, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonValue)
}

// GET /users/[id1]....
func (handler *HttpHandler) GetUsers(w http.ResponseWriter, r *http.Request) {

}
