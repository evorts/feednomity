package feedbacks

import (
	"database/sql/driver"
	"time"
)

type Status string

const (
	StatusNotStarted Status = "not-started"
	StatusDraft      Status = "draft"
	StatusFinal      Status = "final"
)

func (s Status) String() string {
	return string(s)
}

func (s *Status) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return s.String(), nil
}

func (s *Status) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	v, ok := src.(string)
	if ok {
		*s = Status(v)
	}
	return nil
}

type Feedback struct {
	Id                   int64                  `db:"id"`
	DistributionId       int64                  `db:"distribution_id"`
	DistributionTopic    string                 `db:"distribution_topic"`
	DistributionObjectId int64                  `db:"distribution_object_id"`
	RangeStart           *time.Time             `db:"range_start"`
	RangeEnd             *time.Time             `db:"range_end"`
	RespondentId         int64                  `db:"respondent_id"`
	RespondentUsername   string                 `db:"respondent_username"`
	RespondentName       string                 `db:"respondent_name"`
	RespondentEmail      string                 `db:"respondent_email"`
	RespondentGroupId    int64                  `db:"respondent_group_id"`
	RespondentGroupName  string                 `db:"respondent_group_name"`
	RespondentOrgId      int64                  `db:"respondent_org_id"`
	RespondentOrgName    string                 `db:"respondent_org_name"`
	RespondentRole       string                 `db:"respondent_role"`
	RespondentAssignment string                 `db:"respondent_assignment"`
	RecipientId          int64                  `db:"recipient_id"`
	RecipientUsername    string                 `db:"recipient_username"`
	RecipientName        string                 `db:"recipient_name"`
	RecipientEmail       string                 `db:"recipient_email"`
	RecipientGroupId     int64                  `db:"recipient_group_id"`
	RecipientGroupName   string                 `db:"recipient_group_name"`
	RecipientOrgId       int64                  `db:"recipient_org_id"`
	RecipientOrgName     string                 `db:"recipient_org_name"`
	RecipientRole        string                 `db:"recipient_role"`
	RecipientAssignment  string                 `db:"recipient_assignment"`
	LinkId               int64                  `db:"link_id"`
	Hash                 string                 `db:"hash"`
	Status               Status                 `db:"status"`
	Content              map[string]interface{} `db:"content"`
	CreatedAt            *time.Time             `db:"created_at"`
	UpdatedAt            *time.Time             `db:"updated_at"`
}

type Log struct {
	Id         int64                  `db:"id"`
	FeedbackId int                    `db:"feedback_id"`
	Action     string                 `db:"action"`
	Values     map[string]interface{} `db:"values"`
	ValuesPrev map[string]interface{} `db:"values_prev"`
	Notes      string                 `db:"notes"`
	At         *time.Time             `db:"at"`
}
