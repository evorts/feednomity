package handler

import (
	"github.com/evorts/feednomity/domain/assessments"
	"github.com/evorts/feednomity/domain/distribution"
	"github.com/evorts/feednomity/domain/feedbacks"
	"github.com/evorts/feednomity/domain/objects"
	"github.com/evorts/feednomity/domain/users"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/session"
	"github.com/evorts/feednomity/pkg/template"
	"net/http"
	"strings"
	"time"
)

type ItemValue struct {
	Rating int    `json:"rating"`
	Note   string `json:"note"`
}
type SubmissionData struct {
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
}

func Form360(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	view := req.GetContext().Get("view").(template.IManager)
	sm := req.GetContext().Get("sm").(session.IManager)

	datasource := req.GetContext().Get("db").(database.IManager)

	log.Log("forms360_handler", "request received")

	var (
		link     distribution.Link
		err      error
		query    = r.URL.Query()
		linkHash string
	)
	linkHash = query.Get("hash")
	if len(linkHash) < 1 {
		_ = view.Render(w, http.StatusBadRequest, "404.html", map[string]interface{}{
			"PageTitle": "Page Not Found",
		})
		return
	}
	linkDomain := distribution.NewLinksDomain(datasource)
	link, err = linkDomain.FindByHash(req.GetContext().Value(), linkHash)
	if err != nil || !link.Published || link.Disabled {
		_ = view.Render(w, http.StatusBadRequest, "404.html", map[string]interface{}{
			"PageTitle": "Page Not Found",
		})
		return
	}
	//detect if the link has reached maximum visits
	if link.UsageLimit > 0 {
		usageCount := linkDomain.LinkVisitsCountById(req.GetContext().Value(), link.Id)
		if usageCount >= link.UsageLimit {
			_ = view.Render(w, http.StatusBadRequest, "404.html", map[string]interface{}{
				"PageTitle": "Page Not Found",
			})
			return
		}
	}
	sm.Put(r.Context(), "link_hash", linkHash)
	sm.Put(r.Context(), "token", req.GetToken())
	// grab link object information
	var (
		o                     []*objects.Object
		do                    []*distribution.Object
		recipient, respondent *objects.Object
	)
	distDomain := distribution.NewDistributionDomain(datasource)
	do, err = distDomain.FindObjectByIds(req.GetContext().Value(), link.DistributionObjectId)
	if err != nil || len(do) < 1 {
		_ = view.Render(w, http.StatusBadRequest, "404.html", map[string]interface{}{
			"PageTitle": "Page Not Found",
		})
		return
	}
	objectDomain := objects.NewObjectDomain(datasource)
	o, err = objectDomain.FindByIds(req.GetContext().Value(), do[0].RecipientId, do[0].RespondentId)
	if err != nil || len(o) < 2 {
		_ = view.Render(w, http.StatusBadRequest, "404.html", map[string]interface{}{
			"PageTitle": "Page Not Found",
		})
		return
	}
	recipient = o[0]
	respondent = o[1]
	usersDomain := users.NewUserDomain(datasource)
	g, err2 := usersDomain.FindGroupByIds(req.GetContext().Value(), recipient.UserGroupId)
	if err2 != nil {
		_ = view.Render(w, http.StatusBadRequest, "404.html", map[string]interface{}{
			"PageTitle": "Page Not Found",
		})
		return
	}
	_ = linkDomain.RecordLinkVisitor(
		req.GetContext().Value(),
		link, r.Header.Get("User-Agent"),
		map[string]interface{}{
			"referer":    r.Referer(),
			"ip":         r.RemoteAddr,
			"respondent": respondent.Name,
		},
	)

	until := time.Now()
	since := until.Add(-5 * 30 * 24 * time.Hour)
	assessmentsDomain := assessments.NewAssessmentDomain(datasource)
	factors, _ := assessmentsDomain.FindTemplateDataByKey(req.GetContext().Value(), "review360")
	if err = view.InjectData("Csrf", req.GetToken()).Render(w, http.StatusOK, "client_form360.html", map[string]interface{}{
		"PageTitle":    "360 Review Form",
		"RatingsLabel": strings.Join(factors.Ratings.Labels, ","),
		"Seq": func(i int) int {
			return i + 1
		},
		"Assessments": assessments.Item{
			Recipient: assessments.Client{
				Name:         recipient.Name,
				Organization: g[0].Name,
				Role:         recipient.Role,
				Assignment:   recipient.Assignment,
			},
			Respondent: assessments.Client{
				Name:       respondent.Name,
				Role:       respondent.Role,
				Assignment: respondent.Assignment,
			},
			PeriodSince:      &since,
			PeriodUntil:      &until,
			Strengths:        nil,
			NeedImprovements: nil,
			Ratings:          factors.Ratings.Values,
			RatingsLabel:     factors.Ratings.Labels,
			Factors:          &factors.Factors,
		},
	}); err != nil {
		log.Log("form360_handler", err.Error())
	}
}

func Forms(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	view := req.GetContext().Get("view").(template.IManager)
	datasource := req.GetContext().Get("db").(database.IManager)

	log.Log("forms_handler", "request received")

	var (
		link      distribution.Link
		questions []feedbacks.Question
		err       error
		query     = r.URL.Query()
		linkHash  string
	)

	linkHash = query.Get("hash")
	if len(linkHash) < 1 {
		//todo display page error
		return
	}

	linkDomain := distribution.NewLinksDomain(datasource)
	link, err = linkDomain.FindByHash(req.GetContext().Value(), linkHash)
	if err != nil {
		//todo display page error
		return
	}
	substanceDomain := feedbacks.NewSubstanceDomain(datasource)
	questions, err = substanceDomain.FindQuestionsByGroupId(req.GetContext().Value(), int64(link.DistributionObjectId))
	if err != nil {
		//todo display page error
		return
	}
	_ = linkDomain.RecordLinkVisitor(
		req.GetContext().Value(),
		link, r.Header.Get("User-Agent"),
		map[string]interface{}{
			"referer": r.Referer(),
			"ip":      r.RemoteAddr,
		},
	)
	if err = view.InjectData("Csrf", req.GetToken()).Render(w, http.StatusOK, "forms.html", map[string]interface{}{
		"PageTitle": "Anonymous Feedback Submission Page",
		"Data":      questions,
	}); err != nil {
		log.Log("forms_handler", err.Error())
	}
}
