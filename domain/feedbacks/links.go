package feedbacks

import (
	"context"
	"fmt"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/jackc/pgtype"
	"strings"
	"time"
)

type Link struct {
	Id          int64
	Hash        string
	PIN         string
	GroupId     int64
	Disabled    bool
	Published   bool
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

type linksManager struct {
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
	return &linksManager{dbm: dbm}
}

func (l *linksManager) FindByHash(ctx context.Context, hash string) (link Link, err error) {
	q := fmt.Sprintf(`
		SELECT 
			id, hash, pin, group_id, disabled, usage_limit, published, created_at, updated_at, disabled_at, published_at 
		FROM %s
		WHERE hash = $1`, tableLinks)
	err = l.dbm.QueryRowAndBind(ctx, q, []interface{}{hash}, &link)
	return
}

func (l *linksManager) SaveLinks(ctx context.Context, links []Link) error {
	q := fmt.Sprintf(`
		INSERT INTO %s (hash, pin, group_id, disabled, usage_limit, published, created_at, disabled_at, published_at) 
		VALUES`, tableLinks)
	placeholders := make([]string, 0)
	values := make([]interface{}, 0)
	for _, link := range links {
		placeholders = append(placeholders, "(?,?,?,?,?,?,?,?,?)")
		var disabledAt, publishedAt interface{} = nil, nil
		if link.Disabled {
			disabledAt = "NOW()"
		}
		if link.Published {
			publishedAt = "NOW()"
		}
		values = append(values, link.Hash, link.PIN, link.GroupId, link.Disabled, link.UsageLimit, link.Published,
			"NOW()", disabledAt, publishedAt)
	}
	q = fmt.Sprintf(`%s %s`, q, strings.Join(placeholders, ","))
	_, err := l.dbm.Prepare(ctx, "save_links", q)
	if err != nil {
		return err
	}
	cmd, err2 := l.dbm.Exec(ctx, "save_links", values)
	if err2 != nil {
		return err2
	}
	if cmd.RowsAffected() > 0 {
		return nil
	}
	return fmt.Errorf("no rows created")
}

func (l *linksManager) UpdateLink(ctx context.Context, link Link) error {
	if link.Id < 1 {
		return fmt.Errorf("please provide the correct link identifier")
	}
	q := fmt.Sprintf(`
		UPDATE %s 
		SET 
			hash = ?,
			pin = ?,
			group_id = ?,
			disabled = ?,
			usage_limit = ?,
			published = ?,
			updated_at = NOW(),
			disabled_at = ?
			published_at = ?
		WHERE id = ?`, tableLinks)
	_, err := l.dbm.Prepare(ctx, "update_links", q)
	if err != nil {
		return err
	}
	var disabledAt, publishedAt interface{} = nil, nil
	if link.Disabled {
		disabledAt = "NOW()"
	}
	if link.Published {
		publishedAt = "NOW()"
	}
	cmd, err2 := l.dbm.Exec(
		ctx, "update_links",
		link.Hash, link.PIN, link.GroupId, link.Disabled, link.UsageLimit, link.Published,
		disabledAt, publishedAt)
	if err2 != nil {
		return err2
	}
	if cmd.RowsAffected() > 0 {
		return nil
	}
	return fmt.Errorf("no rows created")
}

func (l *linksManager) DisableLinksByIds(ctx context.Context, ids ...int64) (err error) {
	q := fmt.Sprintf(`
		UPDATE %s 
		SET 
			disabled = ?,
			disabled_at = NOW()
		WHERE id = ANY($1)`, tableLinks)
	aids := &pgtype.Int8Array{}
	if err = aids.Set(ids); err != nil {
		return
	}
	cmd, err2 := l.dbm.Exec(ctx, q, aids)
	if err2 != nil {
		return err2
	}
	if cmd.RowsAffected() > 0 {
		return nil
	}
	return fmt.Errorf("no rows updated")
}

func (l *linksManager) RecordLinkVisitor(ctx context.Context, link Link, agent, ref string) error {
	q := fmt.Sprintf(`
		INSERT INTO %s (link_id, at, agent, ref) 
		VALUES($1, NOW(), $2, $3)`, tableLinkVisits)
	cmd, err := l.dbm.Exec(ctx, q, link.Id, agent, ref)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() > 0 {
		return nil
	}
	return fmt.Errorf("no rows recorded")
}
