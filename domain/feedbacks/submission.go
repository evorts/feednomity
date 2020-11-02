package feedbacks

import (
	"context"
	"github.com/evorts/feednomity/pkg/database"
	"time"
)

type MarkAsType string

const (
	MarkedAsFavorite MarkAsType = "favorite"
)

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

type submissionManager struct {
	dbm database.IManager
}

type ISubmission interface {
	SaveSubmission(ctx context.Context, submitted ...Submission) error
	UpdateSubmission(ctx context.Context, submitted Submission) error
	RemoveSubmissionByIds(ctx context.Context, ids ...int64) error
}

func NewSubmissionDomain(dbm database.IManager) ISubmission {
	return &submissionManager{dbm: dbm}
}

func (s *submissionManager) SaveSubmission(ctx context.Context, submitted ...Submission) error {
	panic("implement me")
}

func (s *submissionManager) UpdateSubmission(ctx context.Context, submitted Submission) error {
	panic("implement me")
}

func (s *submissionManager) RemoveSubmissionByIds(ctx context.Context, ids ...int64) error {
	panic("implement me")
}