package handler

import (
	"fmt"
	"github.com/evorts/feednomity/domain/feedbacks"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/session"
	"github.com/evorts/feednomity/pkg/template"
	"github.com/evorts/feednomity/pkg/validate"
	"net/http"
)

func LinksAPI(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	view := req.GetContext().Get("view").(template.IManager)

	log.Log("links_create_api_handler", "request received")

	var payload struct {
		Page  Page  `json:"page"`
		Limit Limit `json:"limit"`
	}

	_ = req.UnmarshallBody(&payload)

	datasource := req.GetContext().Get("db").(database.IManager)
	linkDomain := feedbacks.NewLinksDomain(datasource)

	links, err := linkDomain.FindLinks(req.GetContext().Value(), payload.Page.Value(), payload.Limit.Value())
	if err != nil {
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LNK:ERR:FND",
				Message: "Bad Request! Some problems occurred when searching the data.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	_ = view.RenderJson(w, http.StatusOK, api.Response{
		Status:  http.StatusOK,
		Content: map[string]interface{}{
			"links": links,
		},
		Error:   nil,
	})
}

func LinksCreateAPI(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	view := req.GetContext().Get("view").(template.IManager)

	log.Log("links_create_api_handler", "request received")

	var payload struct {
		Csrf  string           `json:"csrf"`
		Links []feedbacks.Link `json:"links"`
	}
	err := req.UnmarshallBody(&payload)
	if err != nil || len(payload.Links) < 1 {
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LNK:ERR:BND",
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
	for li, link := range payload.Links {
		if len(link.Hash) < 1 {
			errs[fmt.Sprintf("%d_hash", li)] = "invalid hash"
		}
		if link.GroupId < 1 {
			errs[fmt.Sprintf("%d_group_id", li)] = "invalid group"
		}
		if len(link.PIN) > 0 && len(link.PIN) != 6 {
			errs[fmt.Sprintf("%d_pin", li)] = "pin must be 6 character length"
		}
	}
	if len(errs) > 0 {
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LNK:ERR:VAL",
				Message: "Bad Request! Your request resulting validation error.",
				Reasons: errs,
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	datasource := req.GetContext().Get("db").(database.IManager)
	linkDomain := feedbacks.NewLinksDomain(datasource)
	if err = linkDomain.SaveLinks(req.GetContext().Value(), payload.Links); err != nil {
		_ = view.RenderJson(w, http.StatusExpectationFailed, api.Response{
			Status:  http.StatusExpectationFailed,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LNK:ERR:SAV",
				Message: "Fail to save your request. Please check your data and try again.",
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
	return
}

func LinkUpdateAPI(w http.ResponseWriter, r *http.Request) {

}

func LinksRemoveAPI(w http.ResponseWriter, r *http.Request) {

}
