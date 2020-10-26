package handler

import (
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/template"
	"net/http"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	view := req.GetContext().Get("view").(template.IManager)
	log.Log("404_handler", "request received")
	if err := view.Render(w, "404.html", map[string]interface{}{
		"PageTitle": "404 Page Not Found",
	}); err != nil {
		log.Log("404_handler", err.Error())
	}
}
