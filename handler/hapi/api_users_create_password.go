package hapi

import (
	"fmt"
	"github.com/evorts/feednomity/domain/users"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/memory"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/validate"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
)

func ApiCreatePassword(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).PrepareRestful()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.IManager)
	mem := req.GetContext().Get("mem").(memory.IManager)
	datasource := req.GetContext().Get("db").(database.IManager)

	log.Log("create_password_api_handler", "request received")

	var payload struct {
		Pass    string `json:"pass"`
		Confirm string `json:"confirm"`
		Hash    string `json:"hash"`
	}

	err := req.UnmarshallBody(&payload)
	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "CRP:ERR:BND",
				Message: "Bad Request! Something wrong with the payload of your request.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}

	hash := fmt.Sprintf("fp_%s", payload.Hash)
	// Validate request
	errs := make(map[string]string, 0)
	if !validate.ValidPassword(payload.Pass) {
		errs["user"] = "Not a valid password!"
	}
	if payload.Pass != payload.Confirm {
		errs["confirm"] = "Password confirmation are incorrect!"
	}
	if payload.Hash == "" || mem.GetInt64(req.GetContext().Value(), hash, -1) < 1 {
		errs["global"] = "Cannot proceed your request due to limitation!"
	}
	if len(errs) > 0 {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "CRP:ERR:VAL",
				Message: "Bad Request! Validation error.",
				Reasons: errs,
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	userId := mem.GetInt64(req.GetContext().Value(), hash, -1)
	var items []*users.User
	usersDomain := users.NewUserDomain(datasource)
	items, err = usersDomain.FindByIds(req.GetContext().Value(), userId)
	if err != nil || len(items) < 1 {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "CRP:ERR:USR",
				Message: "Bad Request! not registered.",
				Reasons: map[string]string{},
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	// change password
	err = usersDomain.UpdatePasswordById(req.GetContext().Value(), items[0].Id, payload.Pass)
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
	// remove hash
	_ = mem.Delete(req.GetContext().Value(), hash)
	_ = vm.RenderJson(w, http.StatusOK, api.Response{
		Status:  http.StatusOK,
		Content: make(map[string]interface{}, 0),
	})
}
