package hapi

import (
	"fmt"
	"github.com/evorts/feednomity/domain/distribution"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/config"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/view"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"net/http"
	"time"
)

func ApiLinksCreate(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).PrepareRestful()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.IManager)
	cfg := req.GetContext().Get("cfg").(config.IManager)

	log.Log("links_create_api_handler", "request received")

	var payload struct {
		Items                   []*Link `json:"items"`
		DisableAutoGenerateHash bool    `json:"disable_auto_generate_hash"`
	}
	err := req.UnmarshallBody(&payload)
	if err != nil || len(payload.Items) < 1 {
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
	expireAt := time.Now().Add(time.Duration(cfg.GetConfig().App.HashExpire) * time.Hour)
	for li, link := range payload.Items {
		hash := link.Hash
		if !payload.DisableAutoGenerateHash {
			hash = ksuid.New().String()
			link.Hash = hash
		}
		if len(hash) > 0 {
			payload.Items[li].Hash = hash
		}
		if !distribution.Hash(link.Hash).Valid() {
			errs[fmt.Sprintf("%d_hash", li)] = "invalid hash"
		}
		if len(link.PIN) > 0 && !distribution.PIN(link.PIN).Valid() {
			errs[fmt.Sprintf("%d_pin", li)] = "pin must be 6 character length"
		}
	}
	if len(errs) > 0 {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LNK:ERR:VAL",
				Message: "Bad Request! Your request resulting validation error.",
				Reasons: errs,
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	datasource := req.GetContext().Get("db").(database.IManager)
	linkDomain := distribution.NewLinksDomain(datasource)
	linksId := make([]int64, 0)
	if linksId, err = linkDomain.InsertMultiple(
		req.GetContext().Value(),
		transformLinks(
			payload.Items,
			&expireAt,
			[]string{},
		),
	); err != nil {
		_ = vm.RenderJson(w, http.StatusExpectationFailed, api.Response{
			Status:  http.StatusExpectationFailed,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LNK:ERR:SAV",
				Message: "Fail to save your request. Please check your data and try again.",
				Reasons: map[string]string{
					"save_error": errors.Wrap(err, "something wrong with the execution syntax").Error(),
				},
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	_ = vm.RenderJson(w, http.StatusOK, api.Response{
		Status: http.StatusOK,
		Content: map[string]interface{}{
			"ids": linksId,
		},
	})
}
