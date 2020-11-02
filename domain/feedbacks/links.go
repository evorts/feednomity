package feedbacks

import (
	"context"
	"github.com/evorts/feednomity/pkg/database"
	"time"
)

type Link struct {
	Id          int64
	Hash        string
	PIN         string
	GroupId     int64
	Disabled    bool
	UsageLimit  int64
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	DisabledAt  *time.Time
	PublishedAt *time.Time
}

type LinkVisit struct {
	Id     int64
	LinkId int64
	At     *time.Time
	Agent  string
	Ref    map[string]interface{}
}

type feedbackManager struct {
	dbm database.IManager
}

type IFeedbacks interface {
	FindByHash(ctx context.Context, hash string) (Link, error)
	SaveLinks(ctx context.Context, links []Link) error
	UpdateLink(ctx context.Context, link Link) error
	DisableLinksByIds(ctx context.Context, ids ...int64) error
	RecordLinkVisitor(ctx context.Context, link Link, agent, ref string) error
}

const (
	tableLinks      = "links"
	tableLinkVisits = "link_visits"
)

func NewLinksDomain(dbm database.IManager) IFeedbacks {
	return &feedbackManager{dbm: dbm}
}

func (f *feedbackManager) FindByHash(ctx context.Context, hash string) (Link, error) {
	panic("implement me")
}

func (f *feedbackManager) SaveLinks(ctx context.Context, links []Link) error {
	panic("implement me")
}

func (f *feedbackManager) UpdateLink(ctx context.Context, link Link) error {
	panic("implement me")
}

func (f *feedbackManager) DisableLinksByIds(ctx context.Context, ids ...int64) error {
	panic("implement me")
}

func (f *feedbackManager) RecordLinkVisitor(ctx context.Context, link Link, agent, ref string) error {
	panic("implement me")
}
