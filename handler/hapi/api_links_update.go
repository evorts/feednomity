package hapi

import (
	"github.com/evorts/feednomity/domain/distribution"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/config"
	"github.com/evorts/feednomity/pkg/crypt"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/view"
	"github.com/segmentio/ksuid"
	"net/http"
	"time"
)

func ApiLinkUpdate(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).PrepareRestful()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.IManager)
	aes := req.GetContext().Get("aes").(crypt.ICryptAES)
	cfg := req.GetContext().Get("cfg").(config.IManager)

	log.Log("links_update_api_handler", "request received")

	var payload struct {
		RegenerateHash bool              `json:"regenerate_hash"`
		Item           distribution.Link `json:"item"`
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
	if payload.Item.Id < 1 {
		errs["id"] = "not a valid identifier"
	}
	if !payload.RegenerateHash && !distribution.Hash(payload.Item.Hash).Valid() {
		errs["hash"] = "not a valid hash code"
	}
	if !distribution.PIN(payload.Item.PIN).Valid() {
		errs["pin"] = "pin must be 6 character length"
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
	hh := NewHashHelper(aes)
	expireAt := time.Now().Add(time.Duration(cfg.GetConfig().App.HashExpire) * time.Hour)
	if payload.RegenerateHash {
		payload.Item.Hash = hh.Generate(expireAt, ksuid.New().String(), map[string]interface{}{
			"usage_limit":            payload.Item.UsageLimit,
			"pin":                    payload.Item.PIN,
		})
	}
	datasource := req.GetContext().Get("db").(database.IManager)
	linkDomain := distribution.NewLinksDomain(datasource)
	if err = linkDomain.UpdateLink(req.GetContext().Value(), payload.Item); err != nil {
		_ = vm.RenderJson(w, http.StatusExpectationFailed, api.Response{
			Status:  http.StatusExpectationFailed,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LNK:ERR:UPD",
				Message: "Fail to update your request. Please check your data and try again.",
				Reasons: map[string]string{"update_error": err.Error()},
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	_ = vm.RenderJson(w, http.StatusOK, api.Response{
		Status:  http.StatusOK,
		Content: make(map[string]interface{}, 0),
	})
}
