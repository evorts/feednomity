package hapi

import (
	"fmt"
	"github.com/evorts/feednomity/domain/distribution"
	"github.com/evorts/feednomity/domain/feedbacks"
	"github.com/evorts/feednomity/domain/users"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/config"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/utils"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
	"strings"
)

func ApiDistributionBlast(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).PrepareRestful()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.IManager)
	cfg := req.GetContext().Get("cfg").(config.IManager)

	log.Log("links_blast_api_handler", "request received")

	var (
		err     error
		payload struct {
			ForceRedistribution bool    `json:"force_redistribution"`
			Ids                 []int64 `json:"ids"`
			ObjectIds           []int64 `json:"object_ids"`
		}
	)

	err = req.UnmarshallBody(&payload)
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

	if len(payload.Ids) < 1 || len(payload.ObjectIds) < 1 {
		errs["ids"] = "Not a valid items id"
	}

	datasource := req.GetContext().Get("db").(database.IManager)
	distDomain := distribution.NewDistributionDomain(datasource)

	var (
		distributions    []*distribution.Distribution
		distributionsMap = make(map[int64]*distribution.Distribution, 0)
		distIds          = make([]int64, 0)
		objects          []*distribution.Object

		links    []*distribution.Link
		linksMap = make(map[int64]*distribution.Link, 0)
		linksId  = make([]int64, 0)

		usersData []*users.User
		usersMap  = make(map[int64]*users.User, 0)
		userIds   = make([]int64, 0)

		groups    []*users.Group
		groupsMap = make(map[int64]*users.Group, 0)
		groupIds  = make([]int64, 0)

		organizations    []*users.Organization
		organizationsMap = make(map[int64]*users.Organization, 0)
		organizationIds  = make([]int64, 0)
	)

	if len(payload.Ids) > 0 {
		distIds = payload.Ids
		objects, err = distDomain.FindObjectsByDistributionIds(req.GetContext().Value(), payload.Ids...)
	} else {
		objects, err = distDomain.FindObjectByIds(req.GetContext().Value(), payload.ObjectIds...)
	}
	if err != nil || len(objects) < 1 {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "DIST:ERR:FND",
				Message: "Bad Request! Some problems occurred when searching the data.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
objectLoop:
	for _, item := range objects {
		linksId = append(linksId, item.LinkId)
		if _, ok2 := usersMap[item.RecipientId]; !ok2 {
			usersMap[item.RecipientId] = &users.User{}
			userIds = append(userIds, item.RecipientId)
		}
		if _, ok2 := usersMap[item.RespondentId]; !ok2 {
			usersMap[item.RespondentId] = &users.User{}
			userIds = append(userIds, item.RespondentId)
		}
		if len(distIds) < 1 {
			for _, id := range distIds {
				if id == item.DistributionId {
					continue objectLoop
				}
			}
			distIds = append(distIds, item.DistributionId)
		}
	}
	distributions, err = distDomain.FindByIds(req.GetContext().Value(), distIds...)
	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "DIST:ERR:DIST",
				Message: "Bad Request! Some problems occurred when searching the data.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	usersDomain := users.NewUserDomain(datasource)
	usersData, err = usersDomain.FindByIds(req.GetContext().Value(), userIds...)
	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "DIST:ERR:USR",
				Message: "Bad Request! Some problems occurred when searching the data.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	linkDomain := distribution.NewLinksDomain(datasource)
	links, err = linkDomain.FindByIds(req.GetContext().Value(), linksId...)
	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "DIST:ERR:LNK",
				Message: "Bad Request! Some problems occurred when searching the data.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
usersLoop:
	for _, u := range usersData {
		usersMap[u.Id] = u
		if _, ok2 := groupsMap[u.Id]; !ok2 {
			groupsMap[u.Id] = &users.Group{
				Id: u.GroupId,
			}
			for _, gv := range groupIds {
				if gv == u.GroupId {
					continue usersLoop
				}
			}
			groupIds = append(groupIds, u.GroupId)
		}
	}
	groups, err = usersDomain.FindGroupByIds(req.GetContext().Value(), groupIds...)
	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "DIST:ERR:GRP",
				Message: "Bad Request! Some problems occurred when searching the data.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	for gmk, gmv := range groupsMap {
	groupsLoop:
		for _, gv := range groups {
			if gmv.Id != gv.Id {
				continue
			}
			groupsMap[gmk] = gv
			if _, ok2 := organizationsMap[gv.Id]; !ok2 {
				organizationsMap[gv.Id] = &users.Organization{Id: gv.OrgId}
				for _, ov := range organizationIds {
					if ov == gv.OrgId {
						continue groupsLoop
					}
				}
				organizationIds = append(organizationIds, gv.OrgId)
			}
		}
	}
	organizations, err = usersDomain.FindOrganizationByIds(req.GetContext().Value(), organizationIds...)
	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "DIST:ERR:ORG",
				Message: "Bad Request! Some problems occurred when searching the data.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	for _, org := range organizations {
		for omk, omv := range organizationsMap {
			if omv.Id == org.Id {
				organizationsMap[omk] = org
				break
			}
		}
	}
	for _, link := range links {
		linksMap[link.Id] = link
	}
	distIds = make([]int64, 0)
	for _, dist := range distributions {
		// only proceed distributed when forced to
		if !payload.ForceRedistribution && dist.Distributed {
			continue
		}
		distributionsMap[dist.Id] = dist
		distIds = append(distIds, dist.Id)
	}

	queueItems := make([]*distribution.Queue, 0)
	feeds := make([]*feedbacks.Feedback, 0)

	for _, obj := range objects {
		d, ok := distributionsMap[obj.DistributionId]
		if !ok {
			continue
		}
		link := ""
		linkHash := ""
		linkExpiredAt := ""
		if v, ok2 := linksMap[obj.LinkId]; ok2 {
			linkHash = v.Hash
			link = fmt.Sprintf("%s/mbr/link/%s", cfg.GetConfig().App.BaseUrlWeb, v.Hash)
			linkExpiredAt = v.ExpiredAt.Format("02 Jan 2006 15:04:05")
		}
		respondent := &users.User{}
		if v, ok2 := usersMap[obj.RespondentId]; ok2 {
			respondent = v
		}
		recipient := &users.User{}
		if v, ok2 := usersMap[obj.RecipientId]; ok2 {
			recipient = v
		}
		groupRecipient := &users.Group{}
		orgRecipient := &users.Organization{}

		if v, ok2 := groupsMap[obj.RecipientId]; ok2 {
			groupRecipient = v
			if ov, ok3 := organizationsMap[v.Id]; ok3 {
				orgRecipient = ov
			}
		}

		groupRespondent := &users.Group{}
		orgRespondent := &users.Organization{}

		if v, ok2 := groupsMap[obj.RespondentId]; ok2 {
			groupRespondent = v
			if ov, ok3 := organizationsMap[v.Id]; ok3 {
				orgRespondent = ov
			}
		}
		recipientName := utils.IIf(len(recipient.DisplayName) < 1, strings.Title(recipient.Username), recipient.DisplayName)
		subject := fmt.Sprintf("Request Feedback: %s - For: %s", d.Topic, recipientName)
		feeds = append(feeds, &feedbacks.Feedback{
			DistributionId:       obj.DistributionId,
			DistributionTopic:    distributionsMap[obj.DistributionId].Topic,
			DistributionObjectId: obj.Id,
			RangeStart:           distributionsMap[obj.DistributionId].RangeStart,
			RangeEnd:             distributionsMap[obj.DistributionId].RangeEnd,
			RespondentId:         respondent.Id,
			RespondentUsername:   respondent.Username,
			RespondentName:       respondent.DisplayName,
			RespondentEmail:      respondent.Email,
			RespondentGroupId:    groupRespondent.Id,
			RespondentGroupName:  groupRespondent.Name,
			RespondentOrgId:      orgRespondent.Id,
			RespondentOrgName:    orgRespondent.Name,
			RecipientId:          recipient.Id,
			RecipientUsername:    recipient.Username,
			RecipientName:        recipient.DisplayName,
			RecipientEmail:       recipient.Email,
			RecipientGroupId:     groupRecipient.Id,
			RecipientGroupName:   groupRecipient.Name,
			RecipientOrgId:       orgRecipient.Id,
			RecipientOrgName:     orgRecipient.Name,
			LinkId:               obj.LinkId,
			Hash:                 linkHash,
			Status:               feedbacks.StatusNotStarted,
			Content:              make(map[string]interface{}, 0),
		})
		queueItems = append(queueItems, &distribution.Queue{
			DistributionObjectId: obj.Id,
			RecipientId:          obj.RecipientId,
			RespondentId:         obj.RespondentId,
			FromEmail:            cfg.GetConfig().App.Contact.Email,
			ToEmail:              respondent.Email,
			Subject:              subject,
			Template:             cfg.GetConfig().App.ReviewMailTemplate,
			Arguments: map[string]interface{}{
				"subject":              subject,
				"respondent_name":      utils.IIf(len(respondent.DisplayName) < 1, strings.Title(respondent.Username), respondent.DisplayName),
				"recipient_name":       recipientName,
				"recipient_job_role":   recipient.JobRole,
				"recipient_group_name": groupRecipient.Name,
				"recipient_org_name":   orgRecipient.Name,
				"from":                 d.RangeStart.Format("Jan 2006"),
				"to":                   d.RangeEnd.Format("Jan 2006"),
				"link":                 link,
				"expired_at":           linkExpiredAt,
			},
		})
	}
	if len(queueItems) < 1 || len(feeds) < 1 {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "DIST:ERR:QTM",
				Message: "Bad Request! Some problems occurred when searching the data.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	// put into feeds

	// put into queues
	_, err = distDomain.InsertQueues(req.GetContext().Value(), queueItems)
	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "DIST:ERR:QUEUE",
				Message: "Bad Request! Some problems occurred when searching the data.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	// update the distribution status
	_ = distDomain.UpdateStatusAndCountByIds(req.GetContext().Value(), distIds...)
	_ = vm.RenderJson(w, http.StatusOK, api.Response{
		Status:  http.StatusOK,
		Content: map[string]interface{}{},
		Error:   nil,
	})
}
