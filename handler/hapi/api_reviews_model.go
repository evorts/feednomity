package hapi

import (
	"github.com/evorts/feednomity/domain/feedbacks"
	"github.com/evorts/feednomity/pkg/utils"
	"time"
)

type ItemValue struct {
	Rating int    `json:"rating"`
	Note   string `json:"note"`
}

type Error struct {
	Code    string
	Message string
	Err     error
}

type FeedbackRequest struct {
	Id             int64            `json:"id"`
	SubmissionType feedbacks.Status `json:"submission_type"`

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
}

type FeedbackResponse struct {
	Id                   int64                  `json:"id"`
	DistributionId       int64                  `json:"distribution_id"`
	DistributionTopic    string                 `json:"distribution_topic"`
	DistributionObjectId int64                  `json:"distribution_object_id"`
	RangeStart           *time.Time             `json:"range_start"`
	RangeEnd             *time.Time             `json:"range_end"`
	RespondentId         int64                  `json:"respondent_id"`
	RespondentUsername   string                 `json:"respondent_username"`
	RespondentName       string                 `json:"respondent_name"`
	RespondentEmail      string                 `json:"respondent_email"`
	RespondentGroupId    int64                  `json:"respondent_group_id"`
	RespondentGroupName  string                 `json:"respondent_group_name"`
	RespondentOrgId      int64                  `json:"respondent_org_id"`
	RespondentOrgName    string                 `json:"respondent_org_name"`
	RespondentRole       string                 `json:"respondent_role"`
	RespondentAssignment string                 `json:"respondent_assignment"`
	RecipientId          int64                  `json:"recipient_id"`
	RecipientUsername    string                 `json:"recipient_username"`
	RecipientName        string                 `json:"recipient_name"`
	RecipientEmail       string                 `json:"recipient_email"`
	RecipientGroupId     int64                  `json:"recipient_group_id"`
	RecipientGroupName   string                 `json:"recipient_group_name"`
	RecipientOrgId       int64                  `json:"recipient_org_id"`
	RecipientOrgName     string                 `json:"recipient_org_name"`
	RecipientRole        string                 `json:"recipient_role"`
	RecipientAssignment  string                 `json:"recipient_assignment"`
	LinkId               int64                  `json:"link_id"`
	Hash                 string                 `json:"hash"`
	Status               feedbacks.Status       `json:"status"`
	Content              map[string]interface{} `json:"content"`
	CreatedAt            *time.Time             `json:"created_at"`
	UpdatedAt            *time.Time             `json:"updated_at"`
}

func transformFeedbacksReverse(f []*feedbacks.Feedback) (t []*FeedbackResponse) {
	t = make([]*FeedbackResponse, 0)
	for _, fv := range f {
		u := &FeedbackResponse{}
		if err := utils.TransformStruct(u, fv); err == nil {
			t = append(t, u)
		}
	}
	return
}