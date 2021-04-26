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
	"strings"
)

func ApiFeedbackSubmission(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	view := req.GetContext().Get("view").(template.IManager)
	sm := req.GetContext().Get("sm").(session.IManager)
	//hash := req.GetContext().Get("hash").(crypt.ICryptHash)
	datasource := req.GetContext().Get("db").(database.IManager)

	log.Log("feedback_api_handler", "request received")

	var payload struct {
		Csrf         string `json:"csrf"`
		FeedbackHash string `json:"hash"`
		DeviceId     string `json:"device_id"`
		Feedbacks    []struct {
			Id           int64  `json:"id"`
			Sequence     int    `json:"sequence"`
			AnswerEssay  string `json:"answer_essay"`
			AnswerChoice int    `json:"answer_option"`
		} `json:"feedbacks"`
	}
	err := req.UnmarshallBody(&payload)
	if err != nil {
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "FRM:ERR:BND",
				Message: "Bad Request! Something wrong with the payload of your request.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	// validate request
	errs := make(map[string]string, 0)
	if validate.IsEmpty(payload.FeedbackHash) {
		errs["username"] = "Not a valid session!"
	}
	if validate.IsEmpty(payload.DeviceId) {
		errs["password"] = "Not a valid source!"
	}
	// csrf check
	sessionCsrf := sm.Get(r.Context(), "token")
	if validate.IsEmpty(payload.Csrf) || sessionCsrf == nil || payload.Csrf != sessionCsrf.(string) {
		errs["session"] = "Not a valid request session!"
	}
	if len(errs) > 0 {
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "FED:ERR:VAL",
				Message: "Bad Request! Validation error.",
				Reasons: errs,
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	/*var link distribution.Link
	linkDomain := distribution.NewLinksDomain(datasource)
	link, err = linkDomain.FindByHash(req.GetContext().Value(), payload.FeedbackHash)
	if err != nil {
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "FED:ERR:NFL",
				Message: "Bad Request! Data not found.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}*/
	var questions []feedbacks.Question
	substanceDomain := feedbacks.NewSubstanceDomain(datasource)
	questions, err = substanceDomain.FindQuestionsByGroupId(req.GetContext().Value(), 1)
	if err != nil {
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "FED:ERR:NFQ",
				Message: "Bad Request! Related data not found.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	for _, feed := range payload.Feedbacks {
		invalidFeedback := true
		for _, q := range questions {
			if feed.Id != q.Id {
				continue
			}
			if feed.Sequence != q.Sequence {
				continue
			}
			if (q.Expect == feedbacks.QuestionEssay && len(strings.Trim(feed.AnswerEssay, " ")) > 0) ||
				(q.Expect == feedbacks.QuestionMultipleChoice && feed.AnswerChoice > 0) {
				invalidFeedback = false
				break
			}
		}
		if invalidFeedback {
			errs[fmt.Sprintf("%d", feed.Sequence)] = fmt.Sprintf("Invalid answer for question number: %d", feed.Sequence)
		}
	}
	if len(errs) > 0 {
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "FED:ERR:ANS",
				Message: "Bad Request! Validation error: answer is not valid.",
				Reasons: errs,
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	var (
		groups []feedbacks.Group
		group  feedbacks.Group
	)
	groups, err = substanceDomain.FindGroupsByIds(req.GetContext().Value(), 1)
	if err != nil || len(groups) < 1 {
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "FED:ERR:NFG",
				Message: "Bad Request! Not a valid feedback group.",
				Reasons: errs,
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	group = groups[0]
	submittedFeedback := make([]feedbacks.Submission, 0)
	// fill parameter
	for _, feed := range payload.Feedbacks {
		for _, q := range questions {
			if feed.Id != q.Id {
				submittedFeedback = append(submittedFeedback, feedbacks.Submission{
					Hash:           payload.FeedbackHash,
					QuestionId:     q.Id,
					QuestionNumber: feed.Sequence,
					Question:       q.Question,
					GroupId:        q.GroupId,
					GroupTitle:     group.Title,
					InvitationType: group.InvitationType,
					Expect:         q.Expect,
					Options:        q.Options,
					AnswerChoice:   feed.AnswerChoice,
					AnswerEssay:    feed.AnswerEssay,
					MarkedAs:       nil,
				})
				continue
			}
		}
	}
	submissionDomain := feedbacks.NewSubmissionDomain(datasource)
	err = submissionDomain.SaveSubmission(req.GetContext().Value(), submittedFeedback...)
	if err != nil {
		_ = view.RenderJson(w, http.StatusInternalServerError, api.Response{
			Status:  http.StatusInternalServerError,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "FED:ERR:SFL",
				Message: "Unexpected error! Could not save your submission data.",
				Reasons: make(map[string]string, 0),
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
