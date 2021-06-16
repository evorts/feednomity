package hapi

import (
	"encoding/json"
	"github.com/evorts/feednomity/domain/assessments"
	"github.com/evorts/feednomity/domain/feedbacks"
	"github.com/evorts/feednomity/handler/helpers"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/utils"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
	"sort"
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
	var (
		mapByRecipient = make(map[int64]*FeedbackSummaryResponseItem, 0) // map by recipient
		responseItems  = make([]*FeedbackSummaryResponseItem, 0)
		eval           = utils.NewEval()
	)
	for _, feed := range feeds {
		cnt, okc := feed.Content["raw"]
		if !okc || cnt == nil {
			continue
		}
		var cntB []byte
		cntB, err = json.Marshal(cnt)
		if err != nil {
			continue
		}
		var content map[string]interface{}
		err = json.Unmarshal(cntB, &content)
		if err != nil {
			continue
		}
		if _, ok := mapByRecipient[feed.RecipientId]; !ok {
			mapByRecipient[feed.RecipientId] = &FeedbackSummaryResponseItem{
				DistributionId: feed.DistributionId,
				Recipient: Person{
					Id:         feed.RecipientId,
					Name:       feed.RecipientName,
					GroupId:    feed.RecipientGroupId,
					GroupName:  feed.RecipientGroupName,
					OrgId:      feed.RecipientOrgId,
					OrgName:    feed.RecipientOrgName,
					Role:       feed.RecipientRole,
					Assignment: feed.RecipientAssignment,
				},
				TotalScore: 0,
				Rating:     "",
				RangeStart: feed.RangeStart,
				RangeEnd:   feed.RangeEnd,
				Items:      make([]*FeedbackItem, 0),
			}
		}
		assessments.BindToFeedbackFactors("", content, factors.Factors)
		score := helpers.CalculateScore(factors.Factors)
		mapByRecipient[feed.RecipientId].TotalScore = (mapByRecipient[feed.RecipientId].TotalScore + score) /
			utils.IIfF64(mapByRecipient[feed.RecipientId].TotalScore == 0, 1, 2)
		mapByRecipient[feed.RecipientId].Rating = helpers.GetRating(
			eval, factors.Ratings.Labels, factors.Ratings.Threshold,
			mapByRecipient[feed.RecipientId].TotalScore,
		)
		mapByRecipient[feed.RecipientId].Items = append(mapByRecipient[feed.RecipientId].Items, &FeedbackItem{
			Id: feed.Id,
			Respondent: Person{
				Id:         feed.RespondentId,
				Name:       feed.RespondentName,
				GroupId:    feed.RespondentGroupId,
				GroupName:  feed.RespondentGroupName,
				OrgId:      feed.RespondentOrgId,
				OrgName:    feed.RespondentOrgName,
				Role:       feed.RespondentRole,
				Assignment: feed.RespondentAssignment,
			},
			DistributionObjectId: feed.DistributionObjectId,
			Score:                score,
			Rating:               helpers.GetRating(eval, factors.Ratings.Labels, factors.Ratings.Threshold, score),
			Factors:              factors.Factors,
			Strengths:            utils.ArrayInterface(content["strengths"].([]interface{})).ToArrayString().Reduce(),
			NeedImprovements:     utils.ArrayInterface(content["improves"].([]interface{})).ToArrayString().Reduce(),
			Status:               feed.Status,
			UpdatedAt:            feed.UpdatedAt,
		})
	}
	//map to array
	for _, item := range mapByRecipient {
		responseItems = append(responseItems, item)
	}
	sort.Slice(responseItems, func(i, j int) bool {
		return responseItems[i].Recipient.Name < responseItems[j].Recipient.Name
	})
	_ = vm.RenderJson(w, http.StatusOK, api.Response{
		Status: http.StatusOK,
		Content: map[string]interface{}{
			"total": total,
			"items": responseItems,
		},
	})
}
