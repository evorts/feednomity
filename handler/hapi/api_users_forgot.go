package hapi

import (
	"fmt"
	"github.com/evorts/feednomity/domain/users"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/config"
	"github.com/evorts/feednomity/pkg/crypt"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/mailer"
	"github.com/evorts/feednomity/pkg/memory"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/utils"
	"github.com/evorts/feednomity/pkg/validate"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
	"strings"
)

const (
	throttlingEmail = 2 * 60
	hashDefaultExpiration = 5 * 60
)
func ApiForgotPassword(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).PrepareRestful()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.IManager)
	hash := req.GetContext().Get("hash").(crypt.ICryptHash)
	mem := req.GetContext().Get("mem").(memory.IManager)
	mail := req.GetContext().Get("mail").(mailer.IMailer)
	cfg := req.GetContext().Get("cfg").(config.IManager)
	datasource := req.GetContext().Get("db").(database.IManager)

	log.Log("forgot_password_api_handler", "request received")

	var payload struct {
		Username string `json:"username"`
		Hash string `json:"hash"`
	}

	err := req.UnmarshallBody(&payload)
	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "FOG:ERR:BND",
				Message: "Bad Request! Something wrong with the payload of your request.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	// Validate request
	errs := make(map[string]string, 0)
	if !validate.ValidUsername(payload.Username) || !validate.ValidEmail(payload.Username) {
		errs["username"] = "Not a valid username or email!"
	}
	payload.Hash = strings.Trim(payload.Hash, " ")
	if !validate.ValidHash(payload.Hash) {
		errs["global"] = "Invalid request session!"
	}
	if len(errs) > 0 {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "FOG:ERR:VAL",
				Message: "Bad Request! Validation error.",
				Reasons: errs,
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	if mem.GetString(req.GetContext().Value(), payload.Hash, "") != "" {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "FOG:ERR:VAL",
				Message: fmt.Sprintf("Bad Request! You have requested multiple times. Please wait another %d minutes.", throttlingEmail/60),
				Reasons: errs,
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	var user *users.User
	usersDomain := users.NewUserDomain(datasource)
	if validate.ValidEmail(payload.Username) {
		user, err = usersDomain.FindByUserEmail(req.GetContext().Value(), payload.Username)
	} else {
		user, err = usersDomain.FindByUsername(req.GetContext().Value(), payload.Username)
	}
	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "FOG:ERR:USR",
				Message: "Bad Request! not registered.",
				Reasons: map[string]string{"err": err.Error()},
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	// set expiration of request -- 5 minute, throttling
	_ = mem.Set(req.GetContext().Value(), payload.Hash, "fp", 5*60)
	// generate forgot password hash
	fpHash := hash.RenewHash().HashWithoutSalt(fmt.Sprintf("%s:%s", user.Email, payload.Hash))
	// set expiration of fp hash -- 5 minute
	_ = mem.Set(req.GetContext().Value(), fmt.Sprintf("fp_%s", fpHash), user.Id, hashDefaultExpiration)
	// send email forgot password
	content := utils.ReadFile(fmt.Sprintf("%s/%s", cfg.GetConfig().App.MailTemplateDirectory, "forgot-password.html"))
	_, err = mail.SendHtml(
		req.GetContext().Value(),
		[]mailer.Target{{Name: user.DisplayName, Email: user.Email}},
		"Forgot Password Request",
		content,
		map[string]string{
			"requester_name": utils.IIf(len(user.DisplayName) > 0, user.DisplayName, user.Username),
			"recovery_link":  fmt.Sprintf("%s/crp/%s", cfg.GetConfig().App.BaseUrlWeb, fpHash),
		},
	)
	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "FOG:ERR:MAL",
				Message: "Process Error! Failed to continue due to internal error.",
				Reasons: map[string]string{"err": err.Error()},
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	_ = vm.RenderJson(w, http.StatusOK, api.Response{
		Status:  http.StatusOK,
		Content: make(map[string]interface{}, 0),
	})
}
