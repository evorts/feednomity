package hapi

import (
	"github.com/evorts/feednomity/domain/users"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/utils"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
)

func ApiUserUpdate(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).PrepareRestful()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.IManager)

	log.Log("users_update_api_handler", "request received")

	var payload struct {
		Csrf string      `json:"csrf"`
		User *User `json:"user"`
	}

	err := req.UnmarshallBody(&payload)
	if err != nil || payload.User == nil {
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
	user := req.GetUserData()

	if payload.User.Id < 1 {
		errs["id"] = "not a valid identifier"
	}
	if payload.User.GroupId < 1 {
		errs["group_id"] = "not a valid group"
	}
	if len(payload.User.PIN) > 0 && !users.PIN(payload.User.PIN).Valid() {
		errs["pin"] = users.PIN(payload.User.PIN).Rule()
	}
	if len(payload.User.Password) > 0 && !users.PASSWORD(payload.User.Password).Valid() {
		errs["pwd"] = users.PASSWORD(payload.User.Password).Rule()
	}
	// check eligibility of the users to update data
	if len(errs) < 1 && !eligible(
		*user,
		req.GetUserAccessScope(),
		payload.User.Id, payload.User.GroupId,
	) {
		errs["eligibility"] = "Not eligible to make this request."
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
	usersDomain := users.NewUserDomain(datasource)

	var ui []*users.User
	ui, err = usersDomain.FindByIds(req.GetContext().Value(), payload.User.Id)

	if len(ui) < 1 {
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
	var u = ui[0]
	err = utils.MergeStruct(u, payload.User, []string{"Username", "Email", "Phone", "Password", "PIN"})
	if err != nil {
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
	if len(payload.User.Password) < 1 {
		u.Password = ""
	}
	if len(payload.User.PIN) < 1 {
		u.PIN = ""
	}
	if err = usersDomain.Update(req.GetContext().Value(), *u); err != nil {
		_ = vm.RenderJson(w, http.StatusExpectationFailed, api.Response{
			Status:  http.StatusExpectationFailed,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "USR:ERR:UPD",
				Message: "Fail to update your request. Please check your data and try again.",
				Reasons: map[string]string{"update_error": err.Error()},
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
