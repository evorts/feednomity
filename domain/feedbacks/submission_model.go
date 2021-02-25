package feedbacks

import "time"

type Submission struct {
	Id             int64
	Hash           string
	QuestionId     int64
	QuestionNumber int
	Question       string
	GroupId        int64
	GroupTitle     string
	InvitationType InvitationType
	Expect         QuestionType
	Options        []string
	AnswerChoice   int
	AnswerEssay    string
	MarkedAs       []MarkAsType
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
}

type SubmissionAudience struct {
	Id int64
	SubmissionGroupId int64
	AudienceTitle string
	Audiences []string
}

type InvitationType string

const (
	InvitationMultiLink  InvitationType = "multi-link"
	InvitationSingleLink InvitationType = "single-link"
)

type Group struct {
	Id             int64
	Title          string
	InvitationType InvitationType
	Audiences      []int64
	Disabled       bool
	Published      bool
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
	DisabledAt     *time.Time
	PublishedAt    *time.Time
}

type QuestionType string

const (
	QuestionEssay          QuestionType = "essay"
	QuestionMultipleChoice QuestionType = "choice"
)

type Question struct {
	Id         int64
	Sequence   int
	Question   string
	Expect     QuestionType
	Options    []string
	GroupId    int64
	Disabled   bool
	Mandatory  bool
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	DisabledAt *time.Time
}