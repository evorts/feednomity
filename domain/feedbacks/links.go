package feedbacks

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
	"regexp"
	"strings"
	"time"
)

type Hash string

var (
	validHashPattern = regexp.MustCompile("[a-zA-Z0-9]+")
	validPINPattern = regexp.MustCompile("\\d{6}")
)

func (h Hash) Valid() bool {
	return validHashPattern.MatchString(h.Value()) && h.Length() > 0 && h.Length() <= 512
}

func (h Hash) Length() int {
	return len(h)
}

func (h Hash) Value() string {
	return string(h)
}

type PIN string

func (p PIN) Value() string {
	return string(p)
}

func (p PIN) Valid() bool {
	return validPINPattern.MatchString(p.Value())
}

type Link struct {
	Id          int64      `json:"id"`
	Hash        Hash       `json:"hash"`
	PIN         PIN     `json:"pin"`
	GroupId     int64      `json:"group_id"`
	Disabled    bool       `json:"disabled"`
	Published   bool       `json:"published"`
	UsageLimit  int64      `json:"usage_limit"`
	CreatedAt   *time.Time `json:"-"`
	UpdatedAt   *time.Time `json:"-"`
	DisabledAt  *time.Time `json:"-"`
	PublishedAt *time.Time `json:"-"`
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

type ILinks interface {
	FindLinks(ctx context.Context, page, limit int) ([]Link, int, error)
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

func NewLinksDomain(dbm database.IManager) ILinks {
	return &linksManager{dbm: dbm}
}

func (l *linksManager) FindLinks(ctx context.Context, page, limit int) (links []Link, total int, err error) {
	q := fmt.Sprintf(`SELECT count(id) FROM %s`, tableLinks)
	var (
		rows database.Rows
	)
	links = make([]Link, 0)
	err = l.dbm.QueryRowAndBind(ctx, q, nil, &total)
	if err != nil || total < 1 {
		err = errors.Wrap(err, "It looks like the data is not exist")
		return
	}
	q = fmt.Sprintf(`
		SELECT 
			id, hash, pin, group_id, disabled, usage_limit, published, created_at, updated_at, disabled_at, published_at 
		FROM %s ORDER BY id DESC LIMIT %d OFFSET %d`, tableLinks, limit, (page-1)*limit)
	rows, err = l.dbm.Query(ctx, q, nil)
	if err != nil {
		if err == sql.ErrNoRows {
			return links, total, nil
		}
		return
	}
	for rows.Next() {
		var link Link
		err = rows.Scan(
			&link.Id,
			&link.Hash,
			&link.PIN,
			&link.GroupId,
			&link.Disabled,
			&link.UsageLimit,
			&link.Published,
			&link.CreatedAt,
			&link.UpdatedAt,
			&link.DisabledAt,
			&link.PublishedAt,
		)
		if err != nil {
			return
		}
		links = append(links, link)
	}
	return
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
	q = l.dbm.Rebind(ctx, fmt.Sprintf(`%s %s`, q, strings.Join(placeholders, ",")))
	cmd, err2 := l.dbm.Exec(ctx, q, values...)
	if err2 != nil {
		return errors.Wrap(err2, "failed saving links. some errors in constraint or data.")
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
			disabled_at = ?,
			published_at = ?
		WHERE id = ?`, tableLinks)
	q = l.dbm.Rebind(ctx, q)
	var disabledAt, publishedAt interface{} = nil, nil
	if link.Disabled {
		disabledAt = "NOW()"
	}
	if link.Published {
		publishedAt = "NOW()"
	}
	cmd, err2 := l.dbm.Exec(
		ctx, q,
		link.Hash, link.PIN, link.GroupId, link.Disabled, link.UsageLimit, link.Published,
		disabledAt, publishedAt, link.Id)
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
			disabled = true,
			disabled_at = NOW()
		WHERE id = ANY($1)`, tableLinks)
	aids := &pgtype.Int4Array{}
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
