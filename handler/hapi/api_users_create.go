package hapi

import (
	"fmt"
	"github.com/evorts/feednomity/domain/users"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/view"
	"github.com/pkg/errors"
	"net/http"
)

func ApiUserCreate(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).PrepareRestful()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.IManager)

	log.Log("Users_create_api_handler", "request received")

	var payload struct {
		Csrf  string  `json:"csrf"`
		Users []*User `json:"users"`
	}
	err := req.UnmarshallBody(&payload)
	if err != nil || len(payload.Users) < 1 {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "USR:ERR:BND",
				Message: "Bad Request! Something wrong with the payload of your request.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	// let's do validation
	errs := make(map[string]string, 0)
	for usk, usv := range payload.Users {
		if usv.GroupId < 1 {
			errs[fmt.Sprintf("%d_group_id", usk)] = "not a valid group"
		}
		if len(usv.PIN) > 0 && !users.PIN(usv.PIN).Valid() {
			errs[fmt.Sprintf("%d_pin", usk)] = users.PIN(usv.PIN).Rule()
		}
		if len(usv.Password) > 0 && !users.PASSWORD(usv.Password).Valid() {
			errs[fmt.Sprintf("%d_pwd", usk)] = users.PASSWORD(usv.Password).Rule()
		}
	}
	if len(errs) > 0 {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "USR:ERR:VAL",
				Message: "Bad Request! Your request resulting validation error.",
				Reasons: errs,
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	datasource := req.GetContext().Get("db").(database.IManager)
	userDomain := users.NewUserDomain(datasource)
	if err = userDomain.InsertMultiple(req.GetContext().Value(), transformUsers(payload.Users)); err != nil {
		_ = vm.RenderJson(w, http.StatusExpectationFailed, api.Response{
			Status:  http.StatusExpectationFailed,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "USR:ERR:SAV",
				Message: "Fail to save your request. Please check your data and try again.",
				Reasons: map[string]string{
					"save_error": errors.Wrap(err, "something wrong with the execution syntax").Error(),
				},
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
