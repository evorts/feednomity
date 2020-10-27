package handler

import (
	"github.com/evorts/feednomity/domain/users"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/crypt"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/session"
	"github.com/evorts/feednomity/pkg/template"
	"github.com/evorts/feednomity/pkg/validate"
	"net/http"
	"strings"
	"time"
)

func LoginAPI(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	view := req.GetContext().Get("view").(template.IManager)
	sm := req.GetContext().Get("sm").(session.IManager)
	hash := req.GetContext().Get("hash").(crypt.ICrypt)
	datasource := req.GetContext().Get("db").(database.IManager)

	log.Log("login_api_handler", "request received")

	if !req.IsMethodPost() {
		_ = view.RenderJson(w, http.StatusNotAcceptable, api.Response{
			Status:  http.StatusNotAcceptable,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LOG:ERR:MTD",
				Message: "Bad Request! not acceptable.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	if req.IsLoggedIn() {
		_ = view.RenderJson(w, http.StatusContinue, api.Response{
			Status:  http.StatusContinue,
			Content: make(map[string]interface{}, 0),
		})
		return
	}
	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Remember string `json:"remember"`
		Csrf     string `json:"csrf"`
	}
	err := req.UnmarshallBody(&payload)
	if err != nil {
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LOG:ERR:BND",
				Message: "Bad Request! Something wrong with the payload of your request.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	// validate request
	errs := make(map[string]string, 0)
	if !validate.ValidUsername(payload.Username) {
		errs["username"] = "Not a valid username!"
	}
	if !validate.ValidPassword(payload.Password) {
		errs["password"] = "Not a valid password!"
	}
	// csrf check
	sessionCsrf := sm.Get(r.Context(), "token")
	if validate.IsEmpty(payload.Csrf) || sessionCsrf == nil || payload.Csrf != sessionCsrf.(string) {
		errs["session"] = "Not a valid request session!"
	}
	if len(errs) > 0 {
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LOG:ERR:VAL",
				Message: "Bad Request! Validation error.",
				Reasons: errs,
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	var user *users.User
	user, err = users.NewUserDomain(datasource).FindByUsername(req.GetContext().Value(), payload.Username)
	if err != nil {
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LOG:ERR:USR",
				Message: "Bad Request! User not found.",
				Reasons: map[string]string{"err": err.Error()},
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	// ensure the user and password are correct
	passCrypt := hash.Renew().Crypt(payload.Password)
	if strings.ToLower(passCrypt) != strings.ToLower(user.Password) {
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LOG:ERR:ATH",
				Message: "Bad Request! authentication failed.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	rememberExpiration := 3 * 24 * time.Hour
	if len(payload.Remember) > 0 {
		sm.SetSessionLifetime(rememberExpiration)
	}
	if err := sm.RenewToken(req.GetContext().Value()); err != nil {
		_ = view.RenderJson(w, http.StatusFailedDependency, api.Response{
			Status:  http.StatusFailedDependency,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LOG:ERR:SES",
				Message: "Invalid session!",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	sm.Put(req.GetContext().Value(), "user", user)
	//remove token since login success -- client should redirect to respective protected page
	sm.Remove(req.GetContext().Value(), "token")
	_ = view.RenderJson(w, http.StatusOK, api.Response{
		Status:  http.StatusOK,
		Content: make(map[string]interface{}, 0),
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	view := req.GetContext().Get("view").(template.IManager)
	sm := req.GetContext().Get("sm").(session.IManager)

	log.Log("login_handler", "request received")

	if req.IsLoggedIn() {
		http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
		return
	}
	renderData := map[string]interface{}{
		"PageTitle": "Login Page",
	}
	// render login page
	sm.Put(r.Context(), "token", req.GetToken())
	if err := view.InjectData("Csrf", req.GetToken()).Render(w, http.StatusOK, "login.html", renderData); err != nil {
		log.Log("login_handler", err.Error())
	}
}
