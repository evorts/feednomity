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

func ApiDistributionsList(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).PrepareRestful()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.IManager)
	datasource := req.GetContext().Get("db").(database.IManager)

	log.Log("distributions_list_api_handler", "request received")

	var payload struct {
		Page  int    `json:"page"`
		Limit int    `json:"limit"`
	}
	payload.Page = 1
	payload.Limit = 10
	err := req.UnmarshallBody(&payload)
	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest,
			api.NewResponse(
				http.StatusBadRequest, nil,
				api.NewResponseError(
					"OBJ:ERR:BND",
					"Bad Request! Something wrong with the payload of your request.", nil, nil,
				),
			),
		)
		return
	}
	var (
		items []*distribution.Distribution
		total int
	)
	items, total, err = distribution.NewDistributionDomain(datasource).FindAll(req.GetContext().Value(), payload.Page, payload.Limit)
	if err != nil {
		_ = vm.RenderJson(w, http.StatusExpectationFailed,
			api.NewResponse(
				http.StatusExpectationFailed, nil,
				api.NewResponseError(
					"OBJ:ERR:QUE",
					"Failed during inquiry data.", nil, nil,
				),
			),
		)
		return
	}
	_ = vm.RenderJson(w, http.StatusOK, api.Response{
		Status: http.StatusOK,
		Content: map[string]interface{}{
			"total": total,
			"items":  items,
		},
	})
}
