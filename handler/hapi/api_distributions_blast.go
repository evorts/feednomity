package hapi

import (
	"github.com/evorts/feednomity/domain/distribution"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
)

func ApiDistributionBlast(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).PrepareRestful()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.IManager)

	log.Log("links_blast_api_handler", "request received")

	var payload struct {
		Csrf  string  `json:"csrf"`
		Items []int64 `json:"items"`
	}

	err := req.UnmarshallBody(&payload)
	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LNK:ERR:BND",
				Message: "Bad Request! Something wrong with the payload of your request.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	// let's do validation
	errs := make(map[string]string, 0)

	if len(payload.Items) < 1 {
		errs["items"] = "Not a valid items id"
	}

	datasource := req.GetContext().Get("db").(database.IManager)
	distDomain := distribution.NewDistributionDomain(datasource)

	objects, err := distDomain.FindObjectByIds(req.GetContext().Value(), payload.Items...)
	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LNK:ERR:FND",
				Message: "Bad Request! Some problems occurred when searching the data.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	linksId := make([]int64, 0)
	for _, item := range objects {
		linksId = append(linksId, item.LinkId)
	}
	// @todo: doing email blast here
	_ = vm.RenderJson(w, http.StatusOK, api.Response{
		Status: http.StatusOK,
		Content: map[string]interface{}{
			"links": linksId,
		},
		Error: nil,
	})
}
