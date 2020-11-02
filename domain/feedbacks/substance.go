package feedbacks

import (
	"context"
	"github.com/evorts/feednomity/pkg/database"
	"time"
)

type Audience struct {
	Id         int64
	Title      string
	Emails     []string
	Disabled   bool
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	DisabledAt *time.Time
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
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	DisabledAt *time.Time
}

type substanceManager struct {
	dbm database.IManager
}

type ISubstance interface {
	FindAudiencesByIds(ctx context.Context, ids ...int64) ([]Audience, error)
	FindGroupsByIds(ctx context.Context, ids ...int64) ([]Group, error)
	FindQuestionsByGroupId(ctx context.Context, id int64) ([]Question, error)
}

func NewSubstanceDomain(dbm database.IManager) ISubstance {
	return &substanceManager{dbm: dbm}
}

func (s *substanceManager) FindAudiencesByIds(ctx context.Context, ids ...int64) ([]Audience, error) {
	panic("implement me")
}

func (s *substanceManager) FindGroupsByIds(ctx context.Context, ids ...int64) ([]Group, error) {
	panic("implement me")
}

func (s *substanceManager) FindQuestionsByGroupId(ctx context.Context, id int64) ([]Question, error) {
	panic("implement me")
}
