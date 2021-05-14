package hcf

import (
	"context"
	"encoding/json"
	"github.com/evorts/feednomity/domain/assessments"
	"github.com/evorts/feednomity/domain/distribution"
	"github.com/evorts/feednomity/domain/feedbacks"
	"github.com/evorts/feednomity/handler/hapi"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
	"strings"
)

func populateFields(ctx context.Context, lh string, f feedbacks.IFeedback, factors *assessments.Factor) (strengths []string, improvements []string) {
	fd, err := f.FindDetailByHash(ctx, lh)
	if err != nil {
		return
	}
	var (
		content hapi.FeedbackPayload
		cb      []byte
	)
	cb, err = json.Marshal(fd.Content)
	if err != nil {
		return
	}
	err = json.Unmarshal(cb, &content)
	if err != nil {
		return
	}
	var fieldsValue = make(map[string]interface{}, 0)
	content.FilterFieldsAndTransform(nil, fieldsValue)
	if len(fieldsValue) < 1 {
		return
	}
	//todo need refactor -- this is just stupid approach from me
	for k, v := range fieldsValue {
		for _, f1 := range factors.Items {
			if f1.Key == k {
				if vv, ok := v.(hapi.ItemValue); ok {
					f1.Rating = vv.Rating
					f1.Note = vv.Note
				}
				continue
			}
			if len(f1.Items) < 1 {
				continue
			}
			for _, f2 := range f1.Items {
				if f2.Key == k {
					if vv, ok := v.(hapi.ItemValue); ok {
						f2.Rating = vv.Rating
						f2.Note = vv.Note
					}
					continue
				}
				if len(f2.Items) < 1 {
					continue
				}
				for _, f3 := range f2.Items {
					if f3.Key == k {
						if vv, ok := v.(hapi.ItemValue); ok {
							f3.Rating = vv.Rating
							f3.Note = vv.Note
						}
						continue
					}
				}
			}
		}
	}
	if v, ok := fieldsValue["strengths"]; ok {
		if vv, ok2 := v.([]string); ok2 {
			strengths = vv
		}
	}
	if v, ok := fieldsValue["improves"]; ok {
		if vv, ok2 := v.([]string); ok2 {
			improvements = vv
		}
	}
	return strengths, improvements
}

func Review360Form(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.ITemplateManager)
	datasource := req.GetContext().Get("db").(database.IManager)

	log.Log("review_360_handler", "request received")

	var (
		link  distribution.Link
		err   error
		query = r.URL.Query()
	)
	// link hash
	lh := query.Get("hash")

	if len(lh) < 1 {
		_ = vm.Render(w, http.StatusBadRequest, "404.html", map[string]interface{}{
			"PageTitle": "Page Not Found",
		})
		return
	}
	// trace back from hash to respondent data
	link, linkDomain, usageCount, dist, _, recipient, respondent, group, _, er := hapi.QueryAndValidate(req.GetContext().Value(), datasource, req.GetUserData().Id, lh)
	if er != nil {
		_ = vm.Render(w, http.StatusBadRequest, "404.html", map[string]interface{}{
			"PageTitle": "Page Not Found",
		})
		return
	}

	//detect if the link has reached maximum visits
	if link.UsageLimit > 0 {
		if usageCount >= link.UsageLimit {
			_ = vm.Render(w, http.StatusBadRequest, "404.html", map[string]interface{}{
				"PageTitle": "Page Not Found",
			})
			return
		}
	}
	req.GetSession().Put(r.Context(), "link_hash", lh)
	req.GetSession().Put(r.Context(), "token", req.GetToken())

	_ = linkDomain.RecordLinkVisitor(
		req.GetContext().Value(),
		link,
		respondent.Id,
		respondent.Name,
		r.Header.Get("User-Agent"),
		map[string]interface{}{
			"referer":    r.Referer(),
			"ip":         r.RemoteAddr,
			"respondent": respondent.Name,
		},
	)
	assessmentsDomain := assessments.NewAssessmentDomain(datasource)
	factors, _ := assessmentsDomain.FindTemplateDataByKey(req.GetContext().Value(), "review360")

	strengths, improvements := populateFields(req.GetContext().Value(), lh, feedbacks.NewFeedbackDomain(datasource), factors.Factors)
	if strengths == nil {
		strengths = make([]string, factors.StrengthsFieldCount)
	}
	if improvements == nil {
		improvements = make([]string, factors.ImprovementsFieldCount)
	}

	if err = vm.InjectData("Csrf", req.GetToken()).Render(w, http.StatusOK, "member-review-360.html", map[string]interface{}{
		"PageTitle":    factors.Factors.Title,
		"RatingsLabel": strings.Join(factors.Ratings.Labels, ","),
		"Seq": func(i int) int {
			return i + 1
		},
		"Assessments": assessments.Item{
			Recipient: assessments.Client{
				Name:         recipient.Name,
				Organization: group.Name,
				Role:         recipient.Role,
				Assignment:   recipient.Assignment,
			},
			Respondent: assessments.Client{
				Name:       respondent.Name,
				Role:       respondent.Role,
				Assignment: respondent.Assignment,
			},
			PeriodSince:      dist.RangeStart,
			PeriodUntil:      dist.RangeEnd,
			Strengths:        strengths,
			NeedImprovements: improvements,
			Ratings:          factors.Ratings.Values,
			RatingsLabel:     factors.Ratings.Labels,
			Factors:          factors.Factors,
		},
	}); err != nil {
		log.Log("review_360_handler", err.Error())
	}
}

/*func Forms(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.ITemplateManager)
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
	if err = vm.InjectData("Csrf", req.GetToken()).Render(w, http.StatusOK, "forms.html", map[string]interface{}{
		"PageTitle": "Anonymous Feedback Submission Page",
		"Data":      questions,
	}); err != nil {
		log.Log("forms_handler", err.Error())
	}
}*/
