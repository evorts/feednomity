package handler

import (
	"github.com/evorts/feednomity/domain/feedbacks"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/template"
	"net/http"
)

func Form360(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	view := req.GetContext().Get("view").(template.IManager)
	datasource := req.GetContext().Get("db").(database.IManager)

	log.Log("forms360_handler", "request received")

	var (
		link     feedbacks.Link
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
	linkDomain := feedbacks.NewLinksDomain(datasource)
	link, err = linkDomain.FindByHash(req.GetContext().Value(), linkHash)
	if err != nil {
		/*_ = view.Render(w, http.StatusBadRequest, "404.html", map[string]interface{}{
			"PageTitle": "Page Not Found",
		})*/
		//return
	}
	_ = linkDomain.RecordLinkVisitor(
		req.GetContext().Value(),
		link, r.Header.Get("User-Agent"),
		r.Referer(),
	)
	if err = view.InjectData("Csrf", req.GetToken()).Render(w, http.StatusOK, "client_form360.html", map[string]interface{}{
		"PageTitle": "360 Review Form",
		"Ratings":   []int{1, 2, 3, 4, 5},
		"Factors": feedbacks.AssessmentFactor{
			Seq: func(i int) int {
				return i + 1
			},
			Title:  "360 Review Form",
			Weight: 100,
			Items: []feedbacks.AssessmentFactor{
				{
					Key:   "productivity",
					Title: "PRODUCTIVITY",
					Description: `<p><u>Task completion</u><br>Complete tasks on schedule or as promised, based on priorities.</p>
                                    <p>&nbsp;</p>
                                    <p><u>Deliver more than what was expected</u></p>
                                    <p>&nbsp;</p>
                                    <p><u>Actively seek ways to improve product/project</u></p>`,
					Weight: 30,
				},
				{
					Key:   "quality",
					Title: "QUALITY",
					Description: `<p><u>Deliver great quality of code</u><br>including but not limited to reusable
                                        code, unit test, test case and documentation.</p>
                                    <p>&nbsp;</p>
                                    <p><u>Producing less to no bug</u><br>deliver less to no bugs -- as a result of a
                                        excellent and maintainable code, testing or automation test</p>`,
					Weight: 30,
				},
				{
					Key: "dependability",
					Title:  "DEPENDABILITY",
					Weight: 40,
					Items: []feedbacks.AssessmentFactor{
						{
							Key: "leadership",
							Title:  "3A) LEADERSHIP",
							Weight: 18,
							Items: []feedbacks.AssessmentFactor{
								{
									Key:   "ownership",
									Title: "Objective Driven / Sense of Ownership",
									Description: `<p>Sense of product/project ownership and put team/organisation goal above other
                                        interest.</p>`,
									Weight: 5,
								},
								{
									Key:         "adaptability",
									Title:       "Adaptability",
									Description: `<p>Anticipate, adapt, and adjust to changes.</p>`,
									Weight:      3,
								},
								{
									Key:         "prioritization",
									Title:       "Prioritization",
									Description: `<p>Have a good sense of work prioritization.</p>`,
									Weight:      3,
								},
								{
									Key:   "detail-solving",
									Title: "Attention to Details + Analyze + Problem Solving",
									Description: `<p>Have attention to details, and analyze the situation / challenges based on
                                        acquired info / details.</p>`,
									Weight: 4,
								},
								{
									Key:   "independent",
									Title: "Independent",
									Description: `<p>Able to work by themself without help or influence of others, and high work
                                        ethics even without supervision.</p>`,
									Weight: 3,
								},
							},
						},
						{
							Key: "collaboration",
							Title:  "3B) COLLABORATION",
							Weight: 9,
							Items: []feedbacks.AssessmentFactor{
								{
									Key:         "communication",
									Title:       "Communication and Coordination within and across Team",
									Description: `<p>Active in discussion whether within and across team.</p>`,
									Weight:      5,
								},
								{
									Key:   "inspiring",
									Title: "Inspiring Member",
									Description: `<p>Enthusiastically motivates the team with positivity and create productive
                                        environment.</p>`,
									Weight: 4,
								},
							},
						},
						{
							Key: "responsibility",
							Title:  "3C) RESPONSIBILITY AND COMMITMENT",
							Weight: 13,
							Items: []feedbacks.AssessmentFactor{
								{
									Key:         "integrity",
									Title:       "Integrity + Discipline",
									Description: `<p>Own up to your mistake, be responsible, and be respectful to others.</p>`,
									Weight:      6,
								},
								{
									Key:   "openness",
									Title: "Giving + Receiving Feedback",
									Description: `<p>Able to give honest, constructive criticism and receive feedback as a motivation
                                        for self-growth.</p>`,
									Weight: 3,
								},
								{
									Key:         "extra-mile",
									Title:       "Willingness to go Extra Mile",
									Description: `<p>Committed to their work and willingly take extra work when urgency arises.</p>`,
									Weight:      4,
								},
							},
						},
					},
				},
			},
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
		link      feedbacks.Link
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

	linkDomain := feedbacks.NewLinksDomain(datasource)
	link, err = linkDomain.FindByHash(req.GetContext().Value(), linkHash)
	if err != nil {
		//todo display page error
		return
	}
	substanceDomain := feedbacks.NewSubstanceDomain(datasource)
	questions, err = substanceDomain.FindQuestionsByGroupId(req.GetContext().Value(), link.GroupId)
	if err != nil {
		//todo display page error
		return
	}
	_ = linkDomain.RecordLinkVisitor(
		req.GetContext().Value(),
		link, r.Header.Get("User-Agent"),
		r.Referer(),
	)
	if err = view.InjectData("Csrf", req.GetToken()).Render(w, http.StatusOK, "forms.html", map[string]interface{}{
		"PageTitle": "Anonymous Feedback Submission Page",
		"Data":      questions,
	}); err != nil {
		log.Log("forms_handler", err.Error())
	}
}
