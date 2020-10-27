package handler

import (
	"github.com/evorts/feednomity/domain/feedbacks"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/crypt"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/session"
	"github.com/evorts/feednomity/pkg/template"
	"github.com/evorts/feednomity/pkg/validate"
	"net/http"
)

func Forms(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	view := req.GetContext().Get("view").(template.IManager)
	datasource := req.GetContext().Get("db").(database.IManager)

	log.Log("forms_handler", "request received")

	var (
		link      feedbacks.Link
		questions []feedbacks.Question
		err       error
		query     = r.URL.Query()
		linkHash  string
	)

	linkHash = query.Get("hash")
	if len(linkHash) < 1 {
		//todo display page error
		return
	}

	linkDomain := feedbacks.NewLinksDomain(datasource)
	link, err = linkDomain.FindByHash(req.GetContext().Value(), linkHash)
	if err != nil {
		//todo display page error
		return
	}
	substanceDomain := feedbacks.NewSubstanceDomain(datasource)
	questions, err = substanceDomain.FindQuestionsByGroupId(req.GetContext().Value(), link.GroupId)
	if err != nil {
		//todo display page error
		return
	}
	_ = linkDomain.RecordLinkVisitor(
		req.GetContext().Value(),
		link, r.Header.Get("User-Agent"),
		r.Referer(),
	)
	if err = view.InjectData("Csrf", req.GetToken()).Render(w, http.StatusOK, "forms.html", map[string]interface{}{
		"PageTitle": "Anonymous Feedback Submission Page",
		"Data":      questions,
	}); err != nil {
		log.Log("forms_handler", err.Error())
	}
}

func FeedbackSubmissionAPI(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	view := req.GetContext().Get("view").(template.IManager)
	sm := req.GetContext().Get("sm").(session.IManager)
	hash := req.GetContext().Get("hash").(crypt.ICrypt)
	datasource := req.GetContext().Get("db").(database.IManager)

	log.Log("feedback_api_handler", "request received")

	var payload struct {
		Csrf         string `json:"csrf"`
		FeedbackHash string `json:"hash"`
		DeviceId     string `json:"device_id"`
		Feedbacks    []struct {
		}
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
	var link feedbacks.Link
	linkDomain := feedbacks.NewLinksDomain(datasource)
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
	}
	var questions []feedbacks.Question
	substanceDomain := feedbacks.NewSubstanceDomain(datasource)
	questions, err = substanceDomain.FindQuestionsByGroupId(req.GetContext().Value(), link.GroupId)
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
	// todo validate feedbacks with questions on datasource
	for _, feed := range payload.Feedbacks {
		for _, q := range questions {

		}
	}
	// todo save submission
}
