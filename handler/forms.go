package handler

import (
	"github.com/evorts/feednomity/domain/feedbacks"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/template"
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


