package distribution

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/utils"
	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

type Hash string

var (
	validHashPattern = regexp.MustCompile("[a-zA-Z0-9]+")
	validPINPattern  = regexp.MustCompile("\\d{6}")
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

type linksManager struct {
	dbm database.IManager
}

type ILinksManager interface {
	FindAll(ctx context.Context, page, limit int) ([]Link, int, error)
	FindByHash(ctx context.Context, hash string) (Link, error)
	FindByIds(ctx context.Context, ids ...int64) ([]*Link, error)
	InsertMultiple(ctx context.Context, links []*Link) ([]int64, error)
	UpdateLink(ctx context.Context, link Link) error
	DisableLinksByIds(ctx context.Context, ids ...int64) error
	DeleteLinksByIds(ctx context.Context, ids ...int64) error
	LinkVisitsCountById(ctx context.Context, id int64) int
	RecordLinkVisitor(ctx context.Context, link Link, by int64, byName, agent string, ref map[string]interface{}) error
}

const (
	tableLinks      = "links"
	tableLinkVisits = "link_visits"
)

func NewLinksDomain(dbm database.IManager) ILinksManager {
	return &linksManager{dbm: dbm}
}

func (l *linksManager) FindAll(ctx context.Context, page, limit int) (links []Link, total int, err error) {
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
			id, hash, pin, disabled, usage_limit, published, 
			created_by, updated_by, attributes, expired_at,
			created_at, updated_at, disabled_at, published_at 
		FROM %s ORDER BY id DESC LIMIT %d OFFSET %d`, tableLinks, limit, (page-1)*limit)
	rows, err = l.dbm.Query(ctx, q)
	if err != nil {
		if err == sql.ErrNoRows {
			return links, total, nil
		}
		return
	}
	for rows.Next() {
		var (
			link      Link
			pin       sql.NullString
			updatedBy sql.NullInt64
		)
		err = rows.Scan(
			&link.Id,
			&link.Hash,
			&pin,
			&link.Disabled,
			&link.UsageLimit,
			&link.Published,
			&link.CreatedBy,
			&updatedBy,
			&link.Attributes,
			&link.ExpiredAt,
			&link.CreatedAt,
			&link.UpdatedAt,
			&link.DisabledAt,
			&link.PublishedAt,
		)
		link.PIN = pin.String
		link.UpdatedBy = updatedBy.Int64
		if err != nil {
			return
		}
		links = append(links, link)
	}
	return
}

func (l *linksManager) FindByIds(ctx context.Context, ids ...int64) ([]*Link, error) {
	q := fmt.Sprintf(`
		SELECT 
			id, hash, pin, disabled, usage_limit, published, 
			created_by, updated_by, attributes, expired_at,
			created_at, updated_at, disabled_at, published_at  
		FROM %s
		WHERE id IN (%s)`, tableLinks, strings.TrimRight(strings.Repeat("?,", len(ids)), ","))
	links := make([]*Link, 0)
	rows, err := l.dbm.Query(ctx, l.dbm.Rebind(ctx, q), utils.ArrayInt64(ids).ToArrayInterface()...)
	if err != nil {
		if err == sql.ErrNoRows {
			return links, nil
		}
		return nil, err
	}
	for rows.Next() {
		var (
			link      Link
			pin       sql.NullString
			updatedBy sql.NullInt64
		)
		err = rows.Scan(
			&link.Id,
			&link.Hash,
			&pin,
			&link.Disabled,
			&link.UsageLimit,
			&link.Published,
			&link.CreatedBy,
			&updatedBy,
			&link.Attributes,
			&link.ExpiredAt,
			&link.CreatedAt,
			&link.UpdatedAt,
			&link.DisabledAt,
			&link.PublishedAt,
		)
		link.PIN = pin.String
		link.UpdatedBy = updatedBy.Int64
		if err != nil {
			return links, err
		}
		links = append(links, &link)
	}
	return links, nil
}

func (l *linksManager) FindByHash(ctx context.Context, hash string) (link Link, err error) {
	q := fmt.Sprintf(`
		SELECT 
			id, hash, pin, disabled, usage_limit, published, expired_at, created_at, updated_at, disabled_at, published_at 
		FROM %s
		WHERE hash = $1`, tableLinks)
	var pinDb, hashDb sql.NullString
	err = l.dbm.QueryRowAndBind(ctx, q, []interface{}{hash},
		&link.Id, &hashDb, &pinDb, &link.Disabled, &link.UsageLimit, &link.Published, &link.ExpiredAt,
		&link.CreatedAt, &link.UpdatedAt, &link.DisabledAt, &link.PublishedAt,
	)
	link.PIN = pinDb.String
	link.Hash = hashDb.String
	return
}

func (l *linksManager) InsertMultiple(ctx context.Context, links []*Link) ([]int64, error) {
	q := fmt.Sprintf(`
		INSERT INTO %s (
			hash, pin, disabled, usage_limit, published, created_by,
			expired_at, created_at, disabled_at, published_at
		) 
		VALUES`, tableLinks)
	placeholders := make([]string, 0)
	values := make([]interface{}, 0)
	ids := make([]int64, 0)
	for _, link := range links {
		var (
			pinArg                                   = "?"
			disabledAt, publishedAt, pin interface{} = nil, nil, nil
		)
		if link.Disabled {
			disabledAt = "NOW()"
		}
		if link.Published {
			publishedAt = "NOW()"
		}
		if len(link.PIN) > 0 {
			pinArg = "digest(?, 'sha1')"
			pin = link.PIN
		}
		placeholders = append(placeholders, fmt.Sprintf("(?, %s, ?, ?, ?, ?, ?, NOW(), ?, ?)", pinArg))
		values = append(
			values,
			link.Hash, pin, link.Disabled, link.UsageLimit, link.Published,
			link.CreatedBy, link.ExpiredAt, disabledAt, publishedAt,
		)
	}
	q = l.dbm.Rebind(ctx, fmt.Sprintf(`%s %s RETURNING id`, q, strings.Join(placeholders, ",")))
	rows, err2 := l.dbm.Query(ctx, q, values...)
	if err2 != nil {
		return ids, errors.Wrap(err2, "failed saving links. some errors in constraint or data.")
	}
	//get returning ids here
	for rows.Next() {
		var id int64
		if er := rows.Scan(&id); er != nil {
			continue
		}
		ids = append(ids, id)
	}
	if len(ids) > 0 {
		return ids, nil
	}
	return ids, fmt.Errorf("no rows created")
}

func (l *linksManager) UpdateLink(ctx context.Context, link Link) error {
	if link.Id < 1 {
		return fmt.Errorf("please provide the correct link identifier")
	}
	var (
		pinArg = "?"
		expArg = ""

		disabledAt, publishedAt, pin interface{} = nil, nil, nil
	)
	args := []interface{}{link.Hash}
	if len(link.PIN) > 0 {
		pinArg = "digest(?, 'sha1')"
		pin = link.PIN
		args = append(args, pin)
	}
	if link.Disabled {
		disabledAt = "NOW()"
	}
	if link.Published {
		publishedAt = "NOW()"
	}
	args = append(args, link.Disabled, link.UsageLimit, link.Published, disabledAt, publishedAt)
	if link.ExpiredAt != nil {
		expArg = ", expired_at = ?"
		args = append(args, link.ExpiredAt)
	}
	args = append(args, link.Id)
	q := fmt.Sprintf(`
		UPDATE %s 
		SET 
			hash = ?,
			pin = %s,
			disabled = ?,
			usage_limit = ?,
			published = ?,
			updated_at = NOW(),
			disabled_at = ?,
			published_at = ?
			%s
		WHERE id = ?`, tableLinks, pinArg, expArg)
	q = l.dbm.Rebind(ctx, q)
	cmd, err2 := l.dbm.Exec(
		ctx, q, args...)
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

func (l *linksManager) DeleteLinksByIds(ctx context.Context, ids ...int64) error {
	q := l.dbm.Rebind(ctx, fmt.Sprintf(`
			DELETE FROM %s WHERE id IN (%s)
		`, tableLinks, strings.TrimRight(strings.Repeat("?,", len(ids)), ",")))
	rs, err := l.dbm.Exec(
		ctx,
		q,
		utils.ArrayInt64(ids).ToArrayInterface()...,
	)
	if err != nil {
		return err
	}
	if rs.RowsAffected() < 1 {
		return errors.New("not a single record removed")
	}
	return nil
}

func (l *linksManager) RecordLinkVisitor(ctx context.Context, link Link, by int64, byName, agent string, ref map[string]interface{}) error {
	q := fmt.Sprintf(`
		INSERT INTO %s (link_id, by, by_name, at, agent, ref) 
		VALUES($1, NOW(), $2, $3)`, tableLinkVisits)
	cmd, err := l.dbm.Exec(ctx, q, link.Id, by, byName, agent, ref)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() > 0 {
		return nil
	}
	return fmt.Errorf("no rows recorded")
}

func (l *linksManager) LinkVisitsCountById(ctx context.Context, id int64) (total int) {
	q := fmt.Sprintf(`SELECT count(id) FROM %s`, tableLinks)
	_ = l.dbm.QueryRowAndBind(ctx, q, nil, &total)
	return total
}
