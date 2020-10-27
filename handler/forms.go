package handler

import (
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/template"
	"net/http"
)

func Forms(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	view := req.GetContext().Get("view").(template.IManager)
	log.Log("forms_handler", "request received")
	if err := view.Render(w, http.StatusOK, "forms.html", map[string]interface{}{
		"PageTitle": "Anonymous Feedback Submission Page",
	}); err != nil {
		log.Log("forms_handler", err.Error())
	}
}
