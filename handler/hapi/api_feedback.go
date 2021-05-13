package hapi

import (
	"context"
	"github.com/evorts/feednomity/domain/distribution"
	"github.com/evorts/feednomity/domain/feedbacks"
	"github.com/evorts/feednomity/domain/objects"
	"github.com/evorts/feednomity/domain/users"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/utils"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type ItemValue struct {
	Rating int    `json:"rating"`
	Note   string `json:"note"`
}

type FeedbackPayload struct {
	Productivity  ItemValue `json:"productivity"`
	Quality       ItemValue `json:"quality"`
	Dependability struct {
		Leadership struct {
			Adaptability   ItemValue `json:"adaptability"`
			DetailSolving  ItemValue `json:"detail_solving"`
			Independent    ItemValue `json:"independent"`
			Ownership      ItemValue `json:"ownership"`
			Prioritization ItemValue `json:"prioritization"`
		} `json:"leadership"`
		Collaboration struct {
			Communication ItemValue `json:"communication"`
			Inspiring     ItemValue `json:"inspiring"`
		} `json:"collaboration"`
		Responsibility struct {
			ExtraMile ItemValue `json:"extra_mile"`
			Integrity ItemValue `json:"integrity"`
			Openness  ItemValue `json:"openness"`
		} `json:"responsibility"`
	} `json:"dependability"`
	Strengths        []string `json:"strengths"`
	NeedImprovements []string `json:"improves"`
	Csrf             string   `json:"csrf"`
	Hash             string   `json:"hash"`
	SubmissionType   string   `json:"submission_type"`
}

type Error struct {
	Code    string
	Message string
	Err     error
}

func (d *FeedbackPayload) validate() map[string]string {
	if d == nil {
		return map[string]string{"payload": "There's no valid payload submitted!"}
	}
	var fieldsValue = make(map[string]interface{}, 0)
	d.FilterFieldsAndTransform(nil, fieldsValue)
	if len(fieldsValue) < 1 {
		return map[string]string{"payload": "There's no valid payload submitted!"}
	}
	var errs = make(map[string]string, 0)
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

func (d *FeedbackPayload) FilterFieldsAndTransform(value interface{}, rs map[string]interface{}) {
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
		if typeOfs.Field(i).Type.String() != "handler.ItemValue" && v.Field(i).Kind() == reflect.Struct {
			d.FilterFieldsAndTransform(v.Field(i).Interface(), rs)
			continue
		}
		tagName := strings.Split(tag, "\"")[1]
		rs[tagName] = v.Field(i).Interface()
	}
}

func (d *FeedbackPayload) save(
	ctx context.Context, ds database.IManager,
	link distribution.Link,
	dist *distribution.Distribution,
	distObject *distribution.Object,
	recipient *objects.Object,
	respondent *objects.Object,
	group *users.Group,
	user *users.User,
) (err *Error) {
	// its suppose to be impossible to save when the links actually has been disabled
	if link.Disabled {
		err = &Error{
			Code:    "SUB:ERR:LNK0",
			Message: "Not eligible to do this process!",
		}
		return
	}
	feed := feedbacks.NewFeedbackDomain(ds)
	fd, _ := feed.FindDetailByHash(ctx, d.Hash)
	// when existing data found and already final, then no option to update
	if fd != nil && fd.Status == feedbacks.StatusFinal {
		err = &Error{
			Code:    "SUB:ERR:EXT0",
			Message: "Not eligible to do this process!",
		}
		return
	}
	f, _ := feed.FindByDistId(ctx, dist.Id, distObject.Id)
	// when the distribution also has been disabled it is impossible to save further
	if f != nil && f.Disabled {
		err = &Error{
			Code:    "SUB:ERR:EXT1",
			Message: "Not eligible to do this process!",
		}
		return
	}
	fArgs := feedbacks.Feedback{
		DistributionId:       dist.Id,
		DistributionObjectId: distObject.Id,
		DistributionTopic:    dist.Topic,
		UserGroupId:          group.Id,
		UserGroupName:        group.Name,
		UserId:               user.Id,
		UserName:             user.Username,
		UserDisplayName:      user.DisplayName,
	}
	now := time.Now()
	if f != nil {
		fArgs.Disabled = f.Disabled
		fArgs.DisabledAt = f.DisabledAt
		fArgs.UpdateAt = &now
	}
	fdArgs := feedbacks.Detail{
		LinkId:          link.Id,
		Hash:            link.Hash,
		RespondentId:    respondent.Id,
		RespondentName:  respondent.Name,
		RespondentEmail: respondent.Email,
		RecipientId:     recipient.Id,
		RecipientName:   recipient.Name,
		RecipientEmail:  recipient.Email,
		Content:         d,
		Status:          feedbacks.StatusDraft,
	}
	if d.SubmissionType == feedbacks.StatusFinal.String() {
		fdArgs.Status = feedbacks.StatusFinal
	}
	if fd != nil {
		fdArgs.Id = fd.Id
		fdArgs.UpdatedAt = &now
	}
	er := feed.SaveTx(ctx, fArgs, fdArgs)
	if er != nil {
		err = &Error{
			Code:    "SUB:ERR:SAV1",
			Message: "Saving process failed!",
		}
	}
	return
}

func QueryAndValidate(ctx context.Context, ds database.IManager, userId int64, linkHash string) (
	link distribution.Link,
	linkDomain distribution.ILinks,
	linkUsageCount int,
	dist *distribution.Distribution,
	distObject *distribution.Object,
	recipient *objects.Object,
	respondent *objects.Object,
	group *users.Group,
	user *users.User,
	errs *Error,
) {
	var (
		err error
		d   []*distribution.Distribution
		do  []*distribution.Object
		o   []*objects.Object
		g   []*users.Group
		u   []*users.User
	)
	linkDomain = distribution.NewLinksDomain(ds)
	link, err = linkDomain.FindByHash(ctx, linkHash)
	if err != nil || !link.Published || link.Disabled {
		errs = &Error{
			Code:    "SUB:ERR:ENA0",
			Message: "Submission no longer available!",
			Err:     err,
		}
		return
	}
	distDomain := distribution.NewDistributionDomain(ds)
	do, err = distDomain.FindObjectByRespondentAndLinkId(ctx, userId, link.Id)
	if err != nil || len(do) < 1 {
		errs = &Error{
			Code:    "SUB:ERR:DIO4",
			Message: "Could not find respective information about distribution!",
			Err:     err,
		}
		return
	}
	distObject = do[0]
	d, err = distDomain.FindByIds(ctx, distObject.DistributionId)
	if err != nil || len(d) < 1 {
		errs = &Error{
			Code:    "SUB:ERR:DIS0",
			Message: "Could not find respective information about distribution!",
			Err:     err,
		}
		return
	}
	dist = d[0]
	objectDomain := objects.NewObjectDomain(ds)
	o, err = objectDomain.FindByIds(ctx, distObject.RecipientId, distObject.RespondentId)
	if err != nil || len(o) < 2 {
		errs = &Error{
			Code:    "SUB:ERR:OBJ4",
			Message: "Could not find respective information on objects!",
			Err:     err,
		}
		return
	}
	recipient = o[0]
	respondent = o[1]
	usersDomain := users.NewUserDomain(ds)
	g, err = usersDomain.FindGroupByIds(ctx, recipient.UserGroupId)
	if err != nil || len(g) < 1 {
		errs = &Error{
			Code:    "SUB:ERR:USG4",
			Message: "Could not find respective group of objects!",
			Err:     err,
		}
		return
	}
	group = g[0]
	u, err = usersDomain.FindByIds(ctx, dist.CreatedBy)
	if err != nil || len(u) < 1 {
		errs = &Error{
			Code:    "SUB:ERR:USR4",
			Message: "Could not find respective users owner!",
			Err:     err,
		}
		return
	}
	user = u[0]
	return
}

func Api360Submission(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).PrepareRestful()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.IManager)
	datasource := req.GetContext().Get("db").(database.IManager)

	log.Log("360_submit_handler", "request received")

	var payload *FeedbackPayload

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

	errs := payload.validate()

	if len(errs) > 0 {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.NewResponse(
			http.StatusBadRequest, nil,
			api.NewResponseError(
				"SUB:ERR:VAL",
				"Bad Request! Validation error.", errs, nil,
			),
		),
		)
		return
	}

	link, _, _, dist, distObject, recipient, respondent, group, user, er := QueryAndValidate(req.GetContext().Value(), datasource, req.GetUserData().Id, payload.Hash)

	if er != nil {
		_ = vm.RenderJson(
			w, http.StatusBadRequest, api.NewResponse(
				http.StatusBadRequest, nil,
				api.NewResponseError(
					er.Code,
					er.Message, nil, nil,
				),
			),
		)
		return
	}

	er = payload.save(req.GetContext().Value(), datasource, link, dist, distObject, recipient, respondent, group, user)

	if er != nil {
		_ = vm.RenderJson(
			w, http.StatusBadRequest, api.NewResponse(
				http.StatusBadRequest, nil,
				api.NewResponseError(
					er.Code,
					er.Message, nil, nil,
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
