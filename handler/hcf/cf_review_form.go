package hcf

import (
	"encoding/json"
	"fmt"
	"github.com/evorts/feednomity/domain/assessments"
	"github.com/evorts/feednomity/domain/feedbacks"
	"github.com/evorts/feednomity/handler/hapi"
	"github.com/evorts/feednomity/pkg/config"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
	"strconv"
	"strings"
)

func populateFields(data map[string]interface{}, factors *assessments.Factor) (strengths []string, improvements []string) {
	var (
		content hapi.FeedbackRequest
		cb      []byte
		err     error
	)
	cb, err = json.Marshal(data)
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

func ReviewDetail(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.ITemplateManager)
	ds := req.GetContext().Get("db").(database.IManager)
	cfg := req.GetContext().Get("cfg").(config.IManager)

	log.Log("web_review_form_handler", "request received")

	// get the form if from path
	fid, err := strconv.Atoi(req.GetPathLastValue())
	if err != nil {
		_ = vm.Render(w, http.StatusBadRequest, "404.html", map[string]interface{}{
			"PageTitle": "Page Not Found",
		})
		return
	}
	var (
		feeds []*feedbacks.Feedback
		feed  *feedbacks.Feedback
	)
	feedDomain := feedbacks.NewFeedbackDomain(ds)
	feeds, err = feedDomain.FindByIds(req.GetContext().Value(), int64(fid))
	if err != nil || len(feeds) < 1 {
		_ = vm.Render(w, http.StatusBadRequest, "404.html", map[string]interface{}{
			"PageTitle": "Page Not Found",
		})
		return
	}
	feed = feeds[0]
	if !hapi.Eligible(req.GetUserData(), req.GetUserAccessScope(), feed.RespondentId, feed.RespondentGroupId) {
		_ = vm.Render(w, http.StatusBadRequest, "404.html", map[string]interface{}{
			"PageTitle": "Page Not Found",
		})
		return
	}
	assessmentsDomain := assessments.NewAssessmentDomain(ds)
	factors, _ := assessmentsDomain.FindTemplateDataByKey(req.GetContext().Value(), "review360")
	content := make(map[string]interface{})
	if v, ok := feed.Content["raw"]; ok {
		content = v.(map[string]interface{})
	}
	strengths, improvements := populateFields(content, factors.Factors)
	if strengths == nil {
		strengths = make([]string, factors.StrengthsFieldCount)
	}
	if improvements == nil {
		improvements = make([]string, factors.ImprovementsFieldCount)
	}
	if err = vm.InjectData("Csrf", req.GetToken()).Render(w, http.StatusOK, "member-review-detail.html", map[string]interface{}{
		"PageTitle":    factors.Factors.Title,
		"RatingsLabel": strings.Join(factors.Ratings.Labels, ","),
		"Seq": func(i int) int {
			return i + 1
		},
		"ShowButton": feed.Status != feedbacks.StatusFinal,
		"Id": feed.Id,
		"ApiReviewSubmitUrl": fmt.Sprintf("%s/v1/reviews/submit", cfg.GetConfig().App.BaseUrlApi),
		"Assessments": assessments.Item{
			Recipient: assessments.Client{
				Name:         feed.RecipientName,
				Organization: feed.RecipientGroupName,
				Role:         feed.RecipientRole,
				Assignment:   feed.RecipientAssignment,
			},
			Respondent: assessments.Client{
				Name:       feed.RespondentName,
				Organization: feed.RespondentGroupName,
				Role:       feed.RespondentRole,
				Assignment: feed.RespondentAssignment,
			},
			PeriodSince:      feed.RangeStart,
			PeriodUntil:      feed.RangeEnd,
			Strengths:        strengths,
			NeedImprovements: improvements,
			Ratings:          factors.Ratings.Values,
			RatingsLabel:     factors.Ratings.Labels,
			Factors:          factors.Factors,
		},
	}); err != nil {
		log.Log("web_review_form_handler", err.Error())
	}
}
