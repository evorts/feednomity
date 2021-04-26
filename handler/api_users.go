package handler

import (
	"fmt"
	"github.com/evorts/feednomity/domain/users"
	"github.com/evorts/feednomity/pkg/acl"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/session"
	"github.com/evorts/feednomity/pkg/template"
	"github.com/evorts/feednomity/pkg/utils"
	"github.com/evorts/feednomity/pkg/validate"
	"github.com/pkg/errors"
	"net/http"
)

func ApiUsersList(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	view := req.GetContext().Get("view").(template.IManager)
	sm := req.GetContext().Get("sm").(session.IManager)
	datasource := req.GetContext().Get("db").(database.IManager)

	log.Log("login_api_handler", "request received")

	if !req.IsLoggedIn() {
		_ = view.RenderJson(
			w, http.StatusForbidden,
			api.NewResponse(
				http.StatusForbidden, nil,
				api.NewResponseError("OBJ:ERR:RBD", "Bad Request! now allowed.", nil, nil),
			),
		)
		return
	}
	var payload struct {
		Page  int    `json:"page"`
		Limit int    `json:"limit"`
		Csrf  string `json:"csrf"`
	}
	err := req.UnmarshallBody(&payload)
	if err != nil {
		_ = view.RenderJson(w, http.StatusBadRequest,
			api.NewResponse(
				http.StatusBadRequest, nil,
				api.NewResponseError(
					"OBJ:ERR:BND",
					"Bad Request! Something wrong with the payload of your request.", nil, nil,
				),
			),
		)
		return
	}
	// csrf check
	errs := make(map[string]string, 0)
	sessionCsrf := sm.Get(r.Context(), "token")
	if validate.IsEmpty(payload.Csrf) || sessionCsrf == nil || payload.Csrf != sessionCsrf.(string) {
		errs["session"] = "Not a valid request session!"
	}
	if len(errs) > 0 {
		_ = view.RenderJson(w, http.StatusBadRequest,
			api.NewResponse(
				http.StatusBadRequest, nil,
				api.NewResponseError(
					"OBJ:ERR:VAL",
					"Bad Request! Validation error.", nil, nil,
				),
			),
		)
		return
	}
	var (
		ui    []*users.User
		total int
	)
	ui, total, err = users.NewUserDomain(datasource).FindAll(req.GetContext().Value(), payload.Page, payload.Limit)
	if err != nil {
		_ = view.RenderJson(w, http.StatusExpectationFailed,
			api.NewResponse(
				http.StatusExpectationFailed, nil,
				api.NewResponseError(
					"OBJ:ERR:QUE",
					"Failed during inquiry data.", nil, nil,
				),
			),
		)
		return
	}
	_ = view.RenderJson(w, http.StatusOK, api.Response{
		Status: http.StatusOK,
		Content: map[string]interface{}{
			"total":   total,
			"objects": ui,
		},
	})
}

func ApiUserCreate(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	view := req.GetContext().Get("view").(template.IManager)

	log.Log("Users_create_api_handler", "request received")

	var payload struct {
		Csrf  string  `json:"csrf"`
		Users []*User `json:"users"`
	}
	err := req.UnmarshallBody(&payload)
	if err != nil || len(payload.Users) < 1 {
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
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
	// csrf check
	sm := req.GetContext().Get("sm").(session.IManager)
	sessionCsrf := sm.Get(r.Context(), "token")
	if validate.IsEmpty(payload.Csrf) || sessionCsrf == nil || payload.Csrf != sessionCsrf.(string) {
		errs["session"] = "Not a valid request session!"
	}
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
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
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
		_ = view.RenderJson(w, http.StatusExpectationFailed, api.Response{
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
	_ = view.RenderJson(w, http.StatusOK, api.Response{
		Status:  http.StatusOK,
		Content: make(map[string]interface{}, 0),
	})
}

func ApiUserUpdate(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	view := req.GetContext().Get("view").(template.IManager)

	log.Log("users_update_api_handler", "request received")

	var payload struct {
		Csrf string      `json:"csrf"`
		User *users.User `json:"user"`
	}

	err := req.UnmarshallBody(&payload)
	if err != nil || payload.User == nil {
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
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
	sm := req.GetContext().Get("sm").(session.IManager)
	var user reqio.UserSession
	if err = sm.GetJson(req.GetContext().Value(), "user", &user); err != nil {
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "USR:ERR:SES",
				Message: "Bad Request! Something wrong with the way of your request.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	// csrf check
	sessionCsrf := sm.Get(r.Context(), "token")
	if validate.IsEmpty(payload.Csrf) || sessionCsrf == nil || payload.Csrf != sessionCsrf.(string) {
		errs["session"] = "Not a valid request session!"
	}
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
		user,
		acl.AccessScope(sm.GetString(req.GetContext().Value(), "access_scope")),
		payload.User.Id, payload.User.GroupId,
	) {
		errs["eligibility"] = "Not eligible to make this request."
	}
	if len(errs) > 0 {
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
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
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
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
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
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

	if err = usersDomain.Update(req.GetContext().Value(), *payload.User); err != nil {
		_ = view.RenderJson(w, http.StatusExpectationFailed, api.Response{
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
	_ = view.RenderJson(w, http.StatusOK, api.Response{
		Status:  http.StatusOK,
		Content: make(map[string]interface{}, 0),
	})
}

func ApiUsersDelete(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	view := req.GetContext().Get("view").(template.IManager)

	log.Log("users_delete_api_handler", "request received")

	var payload struct {
		Csrf    string  `json:"csrf"`
		UserIds []int64 `json:"user_ids"`
	}

	err := req.UnmarshallBody(&payload)
	if err != nil {
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
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
	// csrf check
	sm := req.GetContext().Get("sm").(session.IManager)
	sessionCsrf := sm.Get(r.Context(), "token")

	if validate.IsEmpty(payload.Csrf) || sessionCsrf == nil || payload.Csrf != sessionCsrf.(string) {
		errs["session"] = "Not a valid request session!"
	}
	if len(payload.UserIds) < 1 {
		errs["id"] = "not a valid identifier"
	}
	if len(errs) > 0 {
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
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
	if err = usersDomain.DisableByIds(req.GetContext().Value(), payload.UserIds); err != nil {
		_ = view.RenderJson(w, http.StatusExpectationFailed, api.Response{
			Status:  http.StatusExpectationFailed,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "USR:ERR:UPD",
				Message: "Fail to update your request. Please check your data and try again.",
				Reasons: map[string]string{"save_error": err.Error()},
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	_ = view.RenderJson(w, http.StatusOK, api.Response{
		Status:  http.StatusOK,
		Content: make(map[string]interface{}, 0),
	})
}
