package hapi

import (
	"github.com/evorts/feednomity/domain/feedbacks"
	"github.com/evorts/feednomity/pkg/acl"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
)

func ApiReviewList(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).PrepareRestful()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.IManager)
	ds := req.GetContext().Get("db").(database.IManager)

	log.Log("api_review_list_handler", "request received")

	var payload struct {
		Page  Page  `json:"page"`
		Limit Limit `json:"limit"`
	}
	_ = req.UnmarshallBody(&payload)

	feedDomain := feedbacks.NewFeedbackDomain(ds)
	var (
		feeds []*feedbacks.Feedback
		total int
		err error
	)
	switch req.GetUserAccessScope() {
	case acl.AccessScopeGlobal:
		feeds, total, err = feedDomain.FindAll(
			req.GetContext().Value(), payload.Page.Value(), payload.Limit.Value(),
		)
	case acl.AccessScopeOrg:
		feeds, total, err = feedDomain.FindByOrgId(
			req.GetContext().Value(), req.GetUserData().OrgId,
			payload.Page.Value(), payload.Limit.Value(),
		)
	case acl.AccessScopeGroup:
		feeds, total, err = feedDomain.FindByGroupId(
			req.GetContext().Value(), req.GetUserData().GroupId,
			payload.Page.Value(), payload.Limit.Value(),
		)
	default:
		feeds, total, err = feedDomain.FindByRespondentId(
			req.GetContext().Value(), req.GetUserData().Id,
			payload.Page.Value(), payload.Limit.Value(),
		)
	}
	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "FDL:ERR:FND",
				Message: "Bad Request! Some problems occurred when searching the data.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	_ = vm.RenderJson(w, http.StatusOK, api.Response{
		Status: http.StatusOK,
		Content: map[string]interface{}{
			"total": total,
			"items": transformFeedbacksReverse(feeds),
		},
		Error: nil,
	})
}
