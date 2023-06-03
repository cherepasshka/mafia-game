package http_mafia

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/jung-kurt/gofpdf"

	"soa.mafia-game/game-server/domain/models/user"
	"soa.mafia-game/game-server/internal/pdf"
	threadpool "soa.mafia-game/game-server/internal/thread_pool"
)

func GetPdf(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	pdfBytes, err := ioutil.ReadFile(fmt.Sprintf("pdf/%s", filename))
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=document.pdf")
	w.Write(pdfBytes)
}

// GET /users/{login}
func (handler *HttpHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	login := chi.URLParam(r, "login")
	user, exists := handler.users.GetUser(login)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonData, err := json.Marshal(user.Profile)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write(jsonData)
	w.Write([]byte("\n"))

	url := fmt.Sprintf("http://%s/pdf/%s.pdf", r.Host, login)
	fmt.Fprintf(w, "PDF document with user information: %s\n", url)

	// todo: transfer into queue
	pool := threadpool.GetThreadPool()
	pool.AddTask(func() {
		pdf, err := pdf.WriteUser(nil, user)
		if err != nil {
			return
		}
		err = pdf.OutputFileAndClose(fmt.Sprintf("./pdf/%s.pdf", login))
		if err != nil {
			log.Printf("Failed to write user %v", err)
			return
		}
	})
}

// PUT /users/{login}
func (handler *HttpHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	login := chi.URLParam(r, "login")
	old_user, exists := handler.users.GetUser(login)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	new_profile := user.Profile{}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(data, &new_profile); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if old_user.Login != new_profile.Login {
		handler.users.ChangeLogin(old_user.Login, new_profile.Login)
	}
	new_user := old_user
	new_user.Profile = new_profile
	handler.users.SetUser(new_user.Login, new_user)
	handler.users.DeleteUser(old_user.Login)
	w.WriteHeader(http.StatusOK)
}

// POST /users/{login}
func (handler *HttpHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	login := chi.URLParam(r, "login")
	old_user, exists := handler.users.GetUser(login)
	if exists {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	new_profile := user.Profile{}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(data, &new_profile); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if login != new_profile.Login {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	new_user := old_user
	new_user.Profile = new_profile
	handler.users.SetUser(new_profile.Login, new_user)
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
	profiles := make([]user.Profile, 0)
	loginsFilename := ""
	for i := range logins {
		user, exists := handler.users.GetUser(logins[i])
		if exists {
			profiles = append(profiles, user.Profile)
			if len(loginsFilename) == 0 {
				loginsFilename = logins[i]
			} else {
				loginsFilename = fmt.Sprintf("%s-%s", loginsFilename, logins[i])
			}
		}
	}
	jsonData, err := json.Marshal(profiles)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	w.Write(jsonData)
	w.Write([]byte("\n"))

	url := fmt.Sprintf("http://%s/pdf/%s.pdf", r.Host, loginsFilename)
	fmt.Fprintf(w, "PDF document with user information: %s\n", url)
	w.WriteHeader(http.StatusOK)

	// todo: transfer into queue
	pool := threadpool.GetThreadPool()
	pool.AddTask(func() {
		var pdfdoc *gofpdf.Fpdf = nil
		var err error
		for i := range logins {
			user, exists := handler.users.GetUser(logins[i])
			if exists {
				pdfdoc, err = pdf.WriteUser(pdfdoc, user)
				if err != nil {
					return
				}
			}
		}
		err = pdfdoc.OutputFileAndClose(fmt.Sprintf("./pdf/%s.pdf", loginsFilename))
		if err != nil {
			log.Printf("Failed to write user %v", err)
			return
		}

	})
}
