package handler

import (
	"github.com/evorts/feednomity/pkg/config"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/template"
	"net/http"
)

func Reload(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	view := req.GetContext().Get("view").(template.IManager)
	cfg := req.GetContext().Get("config").(config.IManager)
	log.Log("reload_handler", "request received")
	if err := cfg.Reload(); err != nil {
		_ = view.RenderRaw(w,  http.StatusBadGateway, "Error reloading")
		return
	}
	_ = view.RenderRaw(w, http.StatusOK, "Reloading done.")
}
