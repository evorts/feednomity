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

func ApiDistObjectsCreate(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).PrepareRestful()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.IManager)
	cfg := req.GetContext().Get("cfg").(config.IManager)

	log.Log("distributions_objects_create_api_handler", "request received")

	var payload struct {
		Items []*DistributionObject `json:"items"`
	}
	err := req.UnmarshallBody(&payload)
	if err != nil || len(payload.Items) < 1 {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "USR:ERR:BND",
				Message: "Bad Request! Something wrong with the payload of your request.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	// let's do validation
	errs := make(map[string]string, 0)
	for k, v := range payload.Items {
		if v.DistributionId < 1 {
			errs[fmt.Sprintf("%d_dist_id", k)] = "not a valid distribution"
		}
		if v.RecipientId < 1 {
			errs[fmt.Sprintf("%d_recipient_id", k)] = "not a valid recipient"
		}
		if v.RespondentId < 1 {
			errs[fmt.Sprintf("%d_respondent_id", k)] = "not a valid respondent"
		}
	}
	if len(errs) > 0 {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "USR:ERR:VAL",
				Message: "Bad Request! Your request resulting validation error.",
				Reasons: errs,
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	datasource := req.GetContext().Get("db").(database.IManager)
	//create necessary links first
	linkDomain := distribution.NewLinksDomain(datasource)
	links := make([]*distribution.Link, 0)
	linksId := make([]int64, 0)
	expireAt := time.Now().Add(time.Duration(cfg.GetConfig().App.HashExpire) * time.Hour)
	for i := 0; i < len(payload.Items); i++ {
		links = append(links, &distribution.Link{
			Hash: ksuid.New().String(),
			UsageLimit:  1,
			CreatedBy: req.GetUserData().Id,
			ExpiredAt: &expireAt,
		})
	}
	linksId, err = linkDomain.InsertMultiple(req.GetContext().Value(), links)
	distDomain := distribution.NewDistributionDomain(datasource)
	if err = distDomain.InsertObjects(
		req.GetContext().Value(),
		transformDistributionObjects(
			req.GetUserData().Id,
			payload.Items,
			linksId,
			[]string{
				"PublishingStatus", "CreatedAt", "UpdatedAt", "PublishedAt",
			},
		),
	); err != nil {
		_ = vm.RenderJson(w, http.StatusExpectationFailed, api.Response{
			Status:  http.StatusExpectationFailed,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "USR:ERR:SAV",
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
		Status:  http.StatusOK,
		Content: make(map[string]interface{}, 0),
	})
}
