package handler

import (
	"github.com/evorts/godash/pkg/logger"
	"github.com/evorts/godash/pkg/reqio"
	"github.com/evorts/godash/pkg/template"
	"net/http"
)

func Forms(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	view := req.GetContext().Get("view").(template.IManager)
	log.Log("forms_handler", "request received")
	if err := view.Render(w, "forms.html", map[string]interface{}{
		"PageAttributes": map[string]interface{}{
			"Title": "Feedback Form Page",
		},
	}); err != nil {
		log.Log("forms_handler", err.Error())
	}
}
