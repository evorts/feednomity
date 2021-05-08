package hapi

import (
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).PrepareRestful()
	vm := req.GetContext().Get("view").(view.IManager)
	_ = vm.RenderJson(w, http.StatusOK, api.NewResponse(http.StatusOK, map[string]interface{}{}, nil))
}