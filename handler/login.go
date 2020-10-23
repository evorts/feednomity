package handler

import (
	"github.com/evorts/godash/pkg/crypt"
	"github.com/evorts/godash/pkg/db"
	"github.com/evorts/godash/pkg/logger"
	"github.com/evorts/godash/pkg/reqio"
	"github.com/evorts/godash/pkg/session"
	"github.com/evorts/godash/pkg/template"
	"net/http"
	"strings"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	view := req.GetContext().Get("view").(template.IManager)
	sm := req.GetContext().Get("sm").(session.IManager)
	hash := req.GetContext().Get("hash").(crypt.ICrypt)
	log.Log("login_handler", "request received")
	if req.IsLoggedIn() {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	renderData := map[string]interface{}{
		"PageAttributes": map[string]interface{}{
			"Title": "Login Page",
		},
	}
	if req.IsMethodGet() {
		// render login page
		sm.Put(r.Context(), "token", req.GetToken())
		if err := view.Render(w, "login.html", renderData); err != nil {
			log.Log("login_handler", err.Error())
		}
		return
	}
	_ = req.ParseForm()
	// validate form
	user := req.GetFormValue("username")
	pass := req.GetFormValue("password")
	remember := req.GetFormValue("remember")
	csrf := req.GetFormValue("csrf")
	validationErrors := make(map[string]string, 0)
	if len(user) < 1 || strings.TrimSpace(user[0]) == "" {
		validationErrors["username"] = "Please fill up your username correctly"
	}
	if len(pass) < 1 || strings.TrimSpace(pass[0]) == "" {
		validationErrors["password"] = "Please fill up your password correctly"
	}
	if len(csrf) < 1 || strings.TrimSpace(csrf[0]) == "" {
		validationErrors["global"] = "Invalid request session"
	}
	if len(validationErrors) > 0 {
		renderData["Errors"] = validationErrors
		sm.Put(r.Context(), "token", req.GetToken())
		if err := view.Render(w, "login.html", renderData); err != nil {
			log.Log("login_handler", err.Error())
		}
		return
	}
	// csrf check
	sessionCsrf := sm.Get(r.Context(), "token")
	if sessionCsrf == nil || csrf[0] != sessionCsrf.(string) {
		validationErrors["global"] = "Invalid request session"
		renderData["Errors"] = validationErrors
		sm.Put(r.Context(), "token", req.GetToken())
		if err := view.Render(w, "login.html", renderData); err != nil {
			log.Log("login_handler", err.Error())
		}
		return
	}
	// ensure the user and password are correct
	var userFound = &db.User{
		Username: "",
		Password: "",
	}
	passCrypt := hash.Renew().Crypt(pass[0])
	if strings.ToLower(passCrypt) != strings.ToLower(userFound.Password) {
		validationErrors["global"] = "Invalid authentication"
		renderData["Errors"] = validationErrors
		sm.Put(r.Context(), "token", req.GetToken())
		if err := view.Render(w, "login.html", renderData); err != nil {
			log.Log("login_handler", err.Error())
		}
		return
	}
	cookieExpiration := 3 * 24 * time.Hour
	if len(remember) > 0 && len(remember[0]) > 0 {
		sm.SetSessionLifetime(cookieExpiration)
	}
	sm.Put(r.Context(), "user", user[0])
	if err := sm.RenewToken(r.Context()); err != nil {
		validationErrors["global"] = "Failed to process"
		renderData["Errors"] = validationErrors
		sm.Put(r.Context(), "token", req.GetToken())
		if err := view.Render(w, "login.html", renderData); err != nil {
			log.Log("login_handler", err.Error())
		}
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}