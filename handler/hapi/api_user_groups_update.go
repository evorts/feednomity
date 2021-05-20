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

func ApiGroupUpdate(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).PrepareRestful()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.IManager)

	log.Log("groups_update_api_handler", "request received")

	var payload struct {
		Group *UserGroupRequest `json:"group"`
	}

	err := req.UnmarshallBody(&payload)
	if err != nil || payload.Group == nil {
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

	if payload.Group.Id < 1 {
		errs["id"] = "not a valid identifier"
	}
	if payload.Group.OrgId < 1 {
		errs["group_id"] = "not a valid org"
	}
	// check eligibility of the users to update data
	if len(errs) < 1 && !Eligible(
		user,
		req.GetUserAccessScope(),
		user.Id, payload.Group.Id,
	) {
		errs["eligibility"] = "Not Eligible to make this request."
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

	var items []*users.Group
	items, err = usersDomain.FindGroupByIds(req.GetContext().Value(), payload.Group.Id)

	if len(items) < 1 {
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
	var item = items[0]
	err = utils.MergeStruct(item, payload.Group, []string{"Name", "OrgId"})
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
	if err = usersDomain.UpdateGroup(req.GetContext().Value(), *item); err != nil {
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
