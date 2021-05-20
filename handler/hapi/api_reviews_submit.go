package hapi

import (
	"encoding/json"
	"github.com/evorts/feednomity/domain/feedbacks"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/utils"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func (d *FeedbackRequest) Validate() map[string]string {
	errs := make(map[string]string, 0)
	if d == nil {
		errs["payload"] = "There's no valid payload submitted!"
		return errs
	}
	if len(errs) > 0 {
		return errs
	}
	var fieldsValue = make(map[string]interface{}, 0)
	d.FilterFieldsAndTransform(nil, fieldsValue)
	if len(fieldsValue) < 1 {
		errs["payload"] = "There's no valid payload submitted!"
		return errs
	}
	for k, v := range fieldsValue {
		vv, ok := v.(ItemValue)
		if ok && vv.Rating < 1 {
			errs[k] = "field should not be empty!"
			continue
		}
		v3, ok2 := v.([]string)
		if ok2 {
			v3 = utils.ArrayString(v3).Reduce()
			if len(v3) > 0 {
				continue
			}
			errs[k] = "field should not be empty!"
			continue
		}
		v4, ok3 := v.(string)
		if ok3 && k == "hash" && len(v4) < 1 {
			errs[k] = "field should not be empty!"
			continue
		}
	}
	return errs
}

func (d *FeedbackRequest) FilterFieldsAndTransform(value interface{}, rs map[string]interface{}) {
	var v reflect.Value
	if value == nil {
		v = reflect.ValueOf(*d)
	} else {
		v = reflect.ValueOf(value)
	}
	typeOfs := v.Type()
	for i := 0; i < v.NumField(); i++ {
		tag := string(typeOfs.Field(i).Tag)
		if tag == "" || strings.Contains(tag, "\"-\"") {
			continue
		}
		if typeOfs.Field(i).Type.String() != "hapi.ItemValue" && v.Field(i).Kind() == reflect.Struct {
			d.FilterFieldsAndTransform(v.Field(i).Interface(), rs)
			continue
		}
		tagName := strings.Split(tag, "\"")[1]
		rs[tagName] = v.Field(i).Interface()
	}
}

func ApiReviewSubmit(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).PrepareRestful()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.IManager)
	ds := req.GetContext().Get("db").(database.IManager)

	log.Log("api_review_submit_handler", "request received")

	var (
		payload *FeedbackRequest
		feedId int
	)

	err := req.UnmarshallBody(&payload)
	if err != nil || payload == nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "FEED:ERR:BND",
				Message: "Bad Request! Something wrong with the payload of your request.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	feedId, err = strconv.Atoi(payload.Id)
	if err != nil || feedId < 1 {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "FEED:ERR:ID",
				Message: "Bad Request! bad identifier.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
	}
	errs := payload.Validate()

	if len(errs) > 0 {
		_ = vm.RenderJson(w, http.StatusBadRequest,
			api.NewResponse(
				http.StatusBadRequest, nil,
				api.NewResponseError(
					"SUB:ERR:VAL",
					"Bad Request! Validation error.", errs, nil,
				),
			),
		)
		return
	}

	var (
		feeds []*feedbacks.Feedback
		feed  *feedbacks.Feedback
	)

	feedDomain := feedbacks.NewFeedbackDomain(ds)
	feeds, err = feedDomain.FindByIds(req.GetContext().Value(), int64(feedId))

	if err != nil || len(feeds) < 1 || feeds[0].Status == feedbacks.StatusFinal {
		_ = vm.RenderJson(w, http.StatusBadRequest,
			api.NewResponse(
				http.StatusBadRequest, nil,
				api.NewResponseError(
					"SUB:ERR:FED",
					"Information not found or no longer Eligible to be modified!", nil, nil,
				),
			),
		)
		return
	}
	feed = feeds[0]
	var cByte []byte
	cByte, err = json.Marshal(payload)
	feed.Content = map[string]interface{}{
		"raw": payload,
		"enc": cByte, //todo: find way to encrypt the feedback
	}
	feed.Status = feedbacks.Status(
		utils.IIf(
			payload.SubmissionType == feedbacks.StatusFinal,
			payload.SubmissionType.String(),
			feedbacks.StatusDraft.String(),
		),
	)
	err = feedDomain.UpdateStatusAndContent(req.GetContext().Value(), feed.Id, feed.Status, feed.Content)
	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest,
			api.NewResponse(
				http.StatusBadRequest, nil,
				api.NewResponseError(
					"SUB:ERR:SAV1",
					"Saving process failed!", nil, nil,
				),
			),
		)
		return
	}
	_ = vm.RenderJson(w, http.StatusOK, api.Response{
		Status:  http.StatusOK,
		Content: make(map[string]interface{}, 0),
	})
}
