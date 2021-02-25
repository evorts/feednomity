package handler

import (
	"github.com/evorts/feednomity/domain/objects"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/session"
	"github.com/evorts/feednomity/pkg/template"
	"github.com/evorts/feednomity/pkg/validate"
	"net/http"
)

func ObjectListAPI(w http.ResponseWriter, r *http.Request) {
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
		o []*objects.Object
		total int
	)
	o, total, err = objects.NewObjectDomain(datasource).FindAll(req.GetContext().Value(), payload.Page, payload.Limit)
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
		Status:  http.StatusOK,
		Content: map[string]interface{}{
			"total": total,
			"objects": o,
		},
	})
}

func ObjectsCreateAPI(w http.ResponseWriter, r *http.Request) {

}

func ObjectUpdateAPI(w http.ResponseWriter, r *http.Request) {

}

func ObjectsRemoveAPI(w http.ResponseWriter, r *http.Request) {

}
