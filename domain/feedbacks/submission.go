package feedbacks

import (
	"context"
	"github.com/evorts/feednomity/pkg/database"
)

type MarkAsType string

const (
	MarkedAsFavorite MarkAsType = "favorite"
)

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