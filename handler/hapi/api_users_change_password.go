package hapi

import (
	"github.com/evorts/feednomity/domain/users"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/crypt"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/validate"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
	"strings"
)

func ApiChangePassword(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).PrepareRestful()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.IManager)
	hash := req.GetContext().Get("hash").(crypt.ICryptHash)
	datasource := req.GetContext().Get("db").(database.IManager)

	log.Log("change_password_api_handler", "request received")

	var payload struct {
		OldPass string `json:"old_pass"`
		NewPass    string `json:"new_pass"`
		Confirm string `json:"confirm"`
	}

	err := req.UnmarshallBody(&payload)
	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "CHP:ERR:BND",
				Message: "Bad Request! Something wrong with the payload of your request.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	// Validate request
	errs := make(map[string]string, 0)
	if !validate.ValidPassword(payload.OldPass) {
		errs["user"] = "Not a valid old password!"
	}
	if !validate.ValidPassword(payload.NewPass) {
		errs["user"] = "Not a valid new password!"
	}
	if payload.NewPass != payload.Confirm {
		errs["confirm"] = "Password confirmation are incorrect!"
	}
	if len(errs) > 0 {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "CHP:ERR:VAL",
				Message: "Bad Request! Validation error.",
				Reasons: errs,
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	var items []*users.User
	usersDomain := users.NewUserDomain(datasource)
	items, err = usersDomain.FindByIds(req.GetContext().Value(), req.GetUserData().Id)
	if err != nil || len(items) < 1 {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "CHP:ERR:USR",
				Message: "Bad Request! not registered.",
				Reasons: map[string]string{},
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	// ensure the user and password are correct
	passCrypt := hash.RenewHash().HashWithoutSalt(payload.OldPass)
	if passCrypt != strings.TrimLeft(items[0].Password, "\\x") {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "CHP:ERR:OFL",
				Message: "Bad Request! incorrect old password.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}

	// change password
	err = usersDomain.UpdatePasswordById(req.GetContext().Value(), items[0].Id, payload.NewPass)
	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "CHP:ERR:MAL",
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
