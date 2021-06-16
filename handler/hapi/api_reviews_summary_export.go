package hapi

import (
	"fmt"
	"github.com/evorts/feednomity/domain/feedbacks"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/utils"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
	"strings"
)

func ApiSummaryReviewsExport(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).PrepareRestful()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.IManager)
	ds := req.GetContext().Get("db").(database.IManager)

	log.Log("api_summary_export_handler", "request received")

	var payload struct {
		FileType string `json:"file_type"`
		DistributionId int64 `json:"distribution_id"`
	}

	err := req.UnmarshallBody(&payload)

	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "SME:ERR:BND",
				Message: "Bad Request! Something wrong with the payload of your request.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}

	//validation
	errs := make(map[string]string, 0)

	if !utils.InArray([]interface{}{"pdf","xls"}, strings.ToLower(payload.FileType)) {
		errs["filetype"] = "Incorrect filetype"
	}
	if payload.DistributionId < 1 {
		errs["distribution"] = "No distribution argument supplied"
	}

	if len(errs) > 0 {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "SME:ERR:VAL",
				Message: "Bad Request! Your request resulting validation error.",
				Reasons: errs,
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	
	var (
		feeds []*feedbacks.Feedback
	)
	feedDomain := feedbacks.NewFeedbackDomain(ds)
	filters := make(map[string]interface{})
	filters["distribution_id"] = payload.DistributionId
	feeds, err = feedDomain.FindByDistId(req.GetContext().Value(), payload.DistributionId)

	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "SME:ERR:FND",
				Message: "Bad Request! Some problems occurred when searching the data.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}

	//@todo: export to filetype
	fmt.Println(feeds)
}
