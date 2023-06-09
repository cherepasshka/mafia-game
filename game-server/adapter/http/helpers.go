package http_mafia

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

func saveImage(r *http.Request, name string) error {
	image, _, err := r.FormFile("image")
	if err != nil {
		return err
	}
	defer image.Close()
	file, err := os.Create(fmt.Sprintf("images/%s", name))
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, image)
	if err != nil {
		return err
	}
	return nil
}

func getUserProfileFromForm(r *http.Request) user.Profile {
	return user.Profile{
		Login:     r.FormValue("login"),
		Email:     r.FormValue("email"),
		Gender:    user.GenderType(r.FormValue("gender")),
		ImageName: fmt.Sprintf("%s.jpg", r.FormValue("login")),
	}
}

// GET /users/{login}
func (handler *HttpHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	login := chi.URLParam(r, "login")
	user, exists := handler.users.GetUser(login)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User %s not found\n", login)
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

	pool := threadpool.GetThreadPool()
	pool.AddTask(func() {
		pdf, err := pdf.WriteUser(nil, user)
		if err != nil {
			return
		}
		err = pdf.OutputFileAndClose(fmt.Sprintf("./pdf/%s.pdf", login))
		if err != nil {
			log.Printf("!Failed to write user: %v", err)
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
		fmt.Fprintf(w, "User %s not found\n", login)
		return
	}
	new_profile := getUserProfileFromForm(r)

	if old_user.Login != new_profile.Login {
		handler.users.ChangeLogin(old_user.Login, new_profile.Login)
	}
	if len(new_profile.ImageName) > 0 {
		err := saveImage(r, new_profile.ImageName)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Failed to save image: %v", err)
			fmt.Fprintf(w, "Server error\n")
			return
		}
	}
	new_user := old_user
	new_user.Profile = new_profile
	handler.users.SetUser(new_user.Login, new_user)
	if new_user.Login != old_user.Login {
		handler.users.DeleteUser(old_user.Login)
	}
	w.WriteHeader(http.StatusOK)
}

// POST /users/{login}
func (handler *HttpHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	login := chi.URLParam(r, "login")
	old_user, exists := handler.users.GetUser(login)
	if exists {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "User already exists\n")
		return
	}
	new_profile := getUserProfileFromForm(r)
	if login != new_profile.Login {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid login\n")
		return
	}
	if len(new_profile.ImageName) > 0 {
		err := saveImage(r, new_profile.ImageName)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Failed to save image: %v", err)
			fmt.Fprintf(w, "Server error\n")
			return
		}
	}
	new_user := old_user
	new_user.Profile = new_profile
	handler.users.SetUser(new_profile.Login, new_user)
	w.WriteHeader(http.StatusCreated)
}

// DELETE /users/{login}
func (handler *HttpHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	login := chi.URLParam(r, "login")
	_, exists := handler.users.GetUser(login)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User %s not found\n", login)
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
			log.Printf("Failed to write users %v", err)
			return
		}
	})
}
