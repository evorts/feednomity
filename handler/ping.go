package handler

import (
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).PrepareRestful()
	vm := req.GetContext().Get("view").(view.IManager)
	_ = vm.RenderRaw(w, http.StatusOK, "OK")
}