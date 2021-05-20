package hapi

import (
	"github.com/evorts/feednomity/domain/feedbacks"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/utils"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
	"strconv"
)

func ApiReviewDetail(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).PrepareRestful()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.IManager)
	ds := req.GetContext().Get("db").(database.IManager)

	log.Log("api_review_detail_handler", "request received")

	//get id from path
	id, err := strconv.Atoi(req.GetPathLastValue())
	if err != nil || id < 1 {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "FDD:ERR:BND",
				Message: "Bad Request! Some problems occurred when searching the data.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}

	feedDomain := feedbacks.NewFeedbackDomain(ds)
	feeds, err := feedDomain.FindByIds(req.GetContext().Value(), int64(id))
	if err != nil || len(feeds) < 1 {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "FDD:ERR:FND",
				Message: "Bad Request! Some problems occurred when searching the data.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	if !Eligible(req.GetUserData(), req.GetUserAccessScope(), feeds[0].RespondentId, feeds[0].RespondentGroupId) {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "FDD:ERR:FBD",
				Message: "Forbidden! you are not allowed to access this resource.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	var detail FeedbackResponse
	_ = utils.TransformStruct(&detail, feeds[0])
	_ = vm.RenderJson(w, http.StatusOK, api.Response{
		Status: http.StatusOK,
		Content: map[string]interface{}{
			"item": detail,
		},
		Error: nil,
	})
}
