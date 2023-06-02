package http_mafia

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"soa.mafia-game/game-server/domain/models/user"
)

// GET /users/{login}
func (handler *HttpHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	login := chi.URLParam(r, "login")
	user, exists := handler.users.GetUser(login)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	jsonValue, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonValue)
}

// PUT /users/{login}
func (handler *HttpHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	login := chi.URLParam(r, "login")
	old_user, exists := handler.users.GetUser(login)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	new_user := user.User{}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(data, &new_user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if old_user.Login != new_user.Login {
		handler.users.ChangeLogin(old_user.Login, new_user.Login)
	}
	handler.users.SetUser(new_user.Login, new_user)
	handler.users.DeleteUser(old_user.Login)
	w.WriteHeader(http.StatusOK)
}

// POST /users/{login}
func (handler *HttpHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	login := chi.URLParam(r, "login")
	_, exists := handler.users.GetUser(login)
	if exists {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	new_user := user.User{}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(data, &new_user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if login != new_user.Login {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	handler.users.SetUser(new_user.Login, new_user)
	w.WriteHeader(http.StatusOK)
}

// DELETE /users/{login}
func (handler *HttpHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	login := chi.URLParam(r, "login")
	_, exists := handler.users.GetUser(login)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	handler.users.DeleteUser(login)
	handler.users.RemovePlayer(login)
	w.WriteHeader(http.StatusOK)
}

// GET /users/?logins={log1},{log2}....
func (handler *HttpHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query()["logins"]) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	logins := strings.Split(r.URL.Query()["logins"][0], ",")
	users := make([]user.User, 0)
	for i := range logins {
		user, exists := handler.users.GetUser(logins[i])
		if exists {
			users = append(users, user)
		}
	}
	jsonData, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	w.Write(jsonData)
	w.Write([]byte("\n"))
	w.WriteHeader(http.StatusOK)
}
