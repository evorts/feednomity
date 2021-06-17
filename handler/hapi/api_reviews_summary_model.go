package hapi

import (
	"github.com/evorts/feednomity/domain/assessments"
	"github.com/evorts/feednomity/domain/feedbacks"
	"sort"
	"time"
)

type Person struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	GroupId    int64  `json:"group_id,omitempty"`
	GroupName  string `json:"group_name,omitempty"`
	OrgId      int64  `json:"org_id,omitempty"`
	OrgName    string `json:"org_name,omitempty"`
	Role       string `json:"role"`
	Assignment string `json:"assignment"`
}

type FeedbackItem struct {
	Id                   int64                  `json:"id"`
	Respondent           Person                 `json:"respondent"`
	DistributionObjectId int64                  `json:"distribution_object_id"`
	Score                float64                `json:"score"`
	Rating               string                 `json:"rating"`
	Factors              *assessments.Factor    `json:"factors"`
	FactorsRaw           map[string]interface{} `json:"-"`
	Strengths            []string               `json:"strengths"`
	NeedImprovements     []string               `json:"need_improvements"`
	Status               feedbacks.Status       `json:"status"`
	UpdatedAt            *time.Time             `json:"updated_at"`
}

type FeedbackSummaryResponseItem struct {
	DistributionId int64           `json:"distribution_id"`
	Recipient      Person          `json:"recipient"`
	RangeStart     *time.Time      `json:"range_start"`
	RangeEnd       *time.Time      `json:"range_end"`
	TotalScore     float64         `json:"total_score"`
	Rating         string          `json:"rating"`
	Items          []*FeedbackItem `json:"items"`
}

/** distribution section **/

type Respondent struct {
	Person
}

type Recipient struct {
	Person
	Respondents []*Respondent `json:"respondents"`
}

type FeedbackSummaryDistributionItem struct {
	DistributionId    int64        `json:"distribution_id"`
	DistributionTopic string       `json:"distribution_topic"`
	Recipients        []*Recipient `json:"recipients"`
}

func transformFeedbackToSummaryDistribution(f []*feedbacks.Feedback) []*FeedbackSummaryDistributionItem {
	rs := make([]*FeedbackSummaryDistributionItem, 0)
	m := make(map[int64]*FeedbackSummaryDistributionItem, 0)
loopFeedbacks:
	for _, v := range f {
		if _, ok := m[v.DistributionId]; !ok {
			m[v.DistributionId] = &FeedbackSummaryDistributionItem{
				DistributionId:    v.DistributionId,
				DistributionTopic: v.DistributionTopic,
				Recipients:        make([]*Recipient, 0),
			}
		}
		for _, vrc := range m[v.DistributionId].Recipients {
			if vrc.Id != v.RecipientId {
				continue
			}
			vrc.Respondents = append(vrc.Respondents, &Respondent{Person{
				Id:         v.RespondentId,
				Name:       v.RespondentName,
				GroupId:    v.RespondentGroupId,
				GroupName:  v.RespondentGroupName,
				OrgId:      v.RespondentOrgId,
				OrgName:    v.RespondentOrgName,
				Role:       v.RespondentRole,
				Assignment: v.RespondentAssignment,
			}})
			sort.Slice(vrc.Respondents, func(i, j int) bool {
				return vrc.Respondents[i].Name < vrc.Respondents[j].Name
			})
			continue loopFeedbacks
		}
		m[v.DistributionId].Recipients = append(m[v.DistributionId].Recipients, &Recipient{
			Person: Person{
				Id:   v.RecipientId,
				Name: v.RecipientName,
			},
			Respondents: make([]*Respondent, 0),
		})
	}
	for _, vv := range m {
		// sort respondent
		rs = append(rs, vv)
	}
	return rs
}
