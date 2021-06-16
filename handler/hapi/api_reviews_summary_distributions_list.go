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

func ApiSummaryDistributionList(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).PrepareRestful()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.IManager)
	ds := req.GetContext().Get("db").(database.IManager)

	log.Log("api_summary_review_handler", "request received")

	var payload struct {
		Page  Page  `json:"page"`
		Limit Limit `json:"limit"`
	}

	err := req.UnmarshallBody(&payload)

	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "SMR:ERR:BND",
				Message: "Bad Request! Something wrong with the payload of your request.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}

	var (
		feeds []*feedbacks.Feedback
		total int
		errs  = make(map[string]string, 0)
	)

	feedDomain := feedbacks.NewFeedbackDomain(ds)

	switch req.GetUserAccessScope() {
	case acl.AccessScopeGlobal:
		feeds, total, err = feedDomain.SummaryByDistribution(
			req.GetContext().Value(),
			payload.Page.Value(),
			payload.Limit.Value(),
			nil,
		)
	case acl.AccessScopeOrg:
		feeds, total, err = feedDomain.SummaryByDistribution(
			req.GetContext().Value(),
			payload.Page.Value(),
			payload.Limit.Value(),
			map[string]interface{}{
				"recipient_org_id": req.GetUserData().OrgId,
			},
		)
	case acl.AccessScopeGroup:
		feeds, total, err = feedDomain.SummaryByDistribution(
			req.GetContext().Value(),
			payload.Page.Value(),
			payload.Limit.Value(),
			map[string]interface{}{
				"recipient_group_id": req.GetUserData().GroupId,
			},
		)
	default:
		errs["global"] = "Your not eligible to see your peers review"
	}

	if err != nil {
		errs["global"] = err.Error()
	}

	if len(errs) > 0 {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "SMR:ERR:VAL",
				Message: "Bad Request! Your request resulting some error.",
				Reasons: errs,
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	_ = vm.RenderJson(w, http.StatusOK, api.Response{
		Status:  http.StatusOK,
		Content: map[string]interface{}{
			"total": total,
			"items": transformFeedbackToSummaryDistribution(feeds),
		},
	})
}
