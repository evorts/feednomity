package feedbacks

import "time"

type Status string

const (
	StatusDraft Status = "draft"
	StatusFinal Status = "final"
)

func (s Status) String() string {
	return string(s)
}

type Feedback struct {
	Id                   int
	DistributionId       int64
	DistributionObjectId int64
	DistributionTopic    string
	UserGroupId          int64
	UserGroupName        string
	UserId               int64
	UserName             string
	UserDisplayName      string
	Disabled             bool
	CreatedAt            *time.Time
	UpdateAt             *time.Time
	DisabledAt           *time.Time
}

type Detail struct {
	Id              int64
	FeedbackId      int64
	LinkId          int64
	Hash            string
	RespondentId    int64
	RespondentName  string
	RespondentEmail string
	RecipientId     int64
	RecipientName   string
	RecipientEmail  string
	Content         interface{}
	Status          Status
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

type Log struct {
	Id int64
	FeedbackId int
	Action string
	Values map[string]interface{}
	ValuesPrev map[string]interface{}
	Notes string
	At *time.Time
}