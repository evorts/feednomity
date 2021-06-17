package hapi

import (
	"github.com/evorts/feednomity/domain/assessments"
	"github.com/evorts/feednomity/domain/feedbacks"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
)

func ApiSummaryReviews(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).PrepareRestful()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.IManager)
	ds := req.GetContext().Get("db").(database.IManager)

	log.Log("api_summary_review_handler", "request received")

	var payload struct {
		RecipientIds   []int64 `json:"recipient_ids"`
		DistributionId int64   `json:"distribution_id"`
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

	//validation
	errs := make(map[string]string, 0)

	if len(payload.RecipientIds) < 1 {
		errs["recipients"] = "No recipients payload found"
	}
	if payload.DistributionId < 1 {
		errs["distribution_id"] = "No distribution information found"
	}

	if len(errs) > 0 {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "SMR:ERR:VAL",
				Message: "Bad Request! Your request resulting validation error.",
				Reasons: errs,
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	var (
		feeds []*feedbacks.Feedback
		total int
	)
	feedDomain := feedbacks.NewFeedbackDomain(ds)
	filters := make(map[string]interface{})
	filters["recipient_id"] = payload.RecipientIds
	filters["distribution_id"] = payload.DistributionId
	feeds, total, err = feedDomain.FindAllWithFilter(req.GetContext().Value(), filters, true)
	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "SMR:ERR:NOF",
				Message: "Bad Request! something wrong with the result.",
				Reasons: errs,
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	var factors *assessments.Template
	assessmentsDomain := assessments.NewAssessmentDomain(ds)
	factors, err = assessmentsDomain.FindTemplateDataByKey(req.GetContext().Value(), "review360")
	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "SMR:ERR:TPL",
				Message: "Internal error. Could not find factors.",
				Reasons: errs,
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	var responseItems []*FeedbackSummaryResponseItem
	responseItems, err = generateReviewSummaryData(feeds, factors)
	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "SMR:ERR:GEN",
				Message: "Internal error. Could not transform data.",
				Reasons: errs,
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	_ = vm.RenderJson(w, http.StatusOK, api.Response{
		Status: http.StatusOK,
		Content: map[string]interface{}{
			"total": total,
			"items": responseItems,
		},
	})
}
