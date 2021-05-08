package hapi

import (
	"github.com/evorts/feednomity/pkg/config"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
)

func ApiReload(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).PrepareRestful()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.IManager)
	cfg := req.GetContext().Get("config").(config.IManager)
	log.Log("reload_handler", "request received")
	if err := cfg.Reload(); err != nil {
		_ = vm.RenderRaw(w,  http.StatusBadGateway, "Error reloading")
		return
	}
	_ = vm.RenderRaw(w, http.StatusOK, "Reloading done.")
}
