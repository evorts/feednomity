package hcf

import (
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.ITemplateManager)
	log.Log("404_handler", "request received")
	if err := vm.Render(w, http.StatusNotFound, "404.html", map[string]interface{}{
		"PageTitle": "404 Page Not Found",
	}); err != nil {
		log.Log("404_handler", err.Error())
	}
}
