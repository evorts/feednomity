package distribution

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/utils"
	"github.com/pkg/errors"
	"strings"
)

type manager struct {
	dbm database.IManager
}

const (
	tableDistribution = "distributions"
	tableDistributionObjects = "distribution_objects"
)

func NewDistributionDomain(dbm database.IManager) IManager {
	return &manager{dbm: dbm}
}

func (m *manager) FindByIds(ctx context.Context, ids ...int64) ([]*Distribution, error) {
	d := make([]*Distribution, 0)
	q := fmt.Sprintf(`
		SELECT 
			id, topic, created_by, range_start, range_end, disabled, archived, distributed, distribution_limit, distribution_count,  
			created_at, updated_at, disabled_at, archived_at, distributed_at 
		FROM %s WHERE id IN (%s)`, tableDistribution, strings.TrimRight(strings.Repeat("?", len(ids)), ","))
	rows, err := m.dbm.Query(ctx, m.dbm.Rebind(ctx, q), utils.ArrayInt64(ids).ToArrayInterface()...)
	if err != nil {
		if err == sql.ErrNoRows {
			return d, nil
		}
		return nil, err
	}
	for rows.Next() {
		var dd Distribution
		err = rows.Scan(
			&dd.Id,
			&dd.Topic,
			&dd.CreatedBy,
			&dd.RangeStart,
			&dd.RangeEnd,
			&dd.Disabled,
			&dd.Archived,
			&dd.Distributed,
			&dd.DistributionLimit,
			&dd.DistributionCount,
			&dd.CreatedAt,
			&dd.UpdatedAt,
			&dd.DisabledAt,
			&dd.ArchivedAt,
			&dd.DistributedAt,
		)
		if err != nil {
			return nil, err
		}
		d = append(d, &dd)
	}
	return d, nil
}

func (m *manager) FindAll(ctx context.Context, page, limit int) (items []*Distribution, total int, err error) {
	var (
		rows database.Rows
	)
	q := fmt.Sprintf(`SELECT count(id) FROM %s`, tableDistribution)
	items = make([]*Distribution, 0)
	err = m.dbm.QueryRowAndBind(ctx, q, nil, &total)
	if err != nil || total < 1 {
		err = errors.Wrap(err, "It looks like the data is not exist")
		return
	}
	rows, err = m.dbm.Query(
		ctx, fmt.Sprintf(
			`SELECT 
						id, topic, distributed, distribution_limit, distribution_count,
						range_start, range_end, created_by, for_group_id, disabled, archived, 
						created_at, updated_at, disabled_at, archived_at, distributed_at
					FROM %s LIMIT %d OFFSET %d`,
			tableDistribution, limit, (page-1)*limit),
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return items, total, nil
		}
		return
	}
	for rows.Next() {
		var (
			item Distribution
		)
		err = rows.Scan(
			&item.Id,
			&item.Topic,
			&item.Distributed,
			&item.DistributionLimit,
			&item.DistributionCount,
			&item.RangeStart,
			&item.RangeEnd,
			&item.CreatedBy,
			&item.ForGroupId,
			&item.Disabled,
			&item.Archived,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.DisabledAt,
			&item.ArchivedAt,
			&item.DistributedAt,
		)
		if err != nil {
			return
		}
		items = append(items, &item)
	}
	return
}

func (m *manager) InsertMultiple(ctx context.Context, items []*Distribution) error {
	q := fmt.Sprintf(`
		INSERT INTO %s (
			topic, disabled, distributed, distribution_limit, distribution_count,
			range_start, range_end, created_by, for_group_id,
			created_at, disabled_at
		) VALUES`, tableDistribution)
	placeholders := make([]string, 0)
	values := make([]interface{}, 0)
	for _, item := range items {
		var disabledAt interface{} = nil
		placeholders = append(
			placeholders,
			`(?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), ?)`,
		)
		if item.Disabled {
			disabledAt = "NOW()"
		}
		values = append(
			values, item.Topic, item.Disabled, item.Distributed, item.DistributionLimit, item.DistributionCount,
			item.RangeStart, item.RangeEnd, item.CreatedBy, item.ForGroupId, disabledAt,
		)
	}
	q = m.dbm.Rebind(ctx, fmt.Sprintf(`%s %s`, q, strings.Join(placeholders, ",")))
	cmd, err2 := m.dbm.Exec(ctx, q, values...)
	if err2 != nil {
		return errors.Wrap(err2, "failed saving groups. some errors in constraint or data.")
	}
	if cmd.RowsAffected() > 0 {
		return nil
	}
	return fmt.Errorf("no rows created")
}

func (m *manager) Update(ctx context.Context, item Distribution) error {
	if item.Id < 1 {
		return fmt.Errorf("please provide the correct identifier")
	}
	args := []interface{}{
		item.Topic,
		item.Disabled,
		item.Archived,
		item.Distributed,
		item.DistributionLimit,
		item.DistributionCount,
		item.RangeStart,
		item.RangeEnd,
		item.CreatedBy,
		item.ForGroupId,
	}
	var (
		disabledAt interface{} = item.DisabledAt
		archivedAt interface{} = item.ArchivedAt
		distributedAt interface{} = item.DistributedAt
	)
	if item.Disabled {
		disabledAt = "NOW()"
	}
	if item.Archived {
		archivedAt = "NOW()"
	}
	if item.Distributed {
		distributedAt = "NOW()"
	}
	args = append(args, disabledAt)
	args = append(args, archivedAt)
	args = append(args, distributedAt)
	args = append(args, item.Id)
	q := fmt.Sprintf(`
		UPDATE %s 
		SET 
			topic = ?,
			disabled = ?,
			archived = ?,
			distributed = ?,
			distribution_limit = ?,
			distribution_count = ?,
			range_start = ?,
			range_end = ?,
			created_by = ?,
			for_group_id = ?,
			updated_at = NOW(),
			disabled_at = ?,
			archived_at = ?,
			distributed_at = ?
		WHERE id = ?`, tableDistribution)
	q = m.dbm.Rebind(ctx, q)
	cmd, err2 := m.dbm.Exec(ctx, q, args...)
	if err2 != nil {
		return err2
	}
	if cmd.RowsAffected() > 0 {
		return nil
	}
	return fmt.Errorf("no rows updated")
}

func (m *manager) DeleteByIds(ctx context.Context, ids ...int64) error {
	q := m.dbm.Rebind(ctx, fmt.Sprintf(`
			DELETE FROM %s WHERE id IN (%s)
		`, tableDistribution, strings.TrimRight(strings.Repeat("?,", len(ids)), ",")))
	rs, err := m.dbm.Exec(
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

func (m *manager) FindObjectByIds(ctx context.Context, ids ...int64) ([]*Object, error) {
	o := make([]*Object, 0)
	q := fmt.Sprintf(`
		SELECT 
			id, distribution_id, recipient_id, respondent_id, publishing_status, publishing_log, 
			created_at, updated_at, published_at 
		FROM %s WHERE id IN (%s)`, tableDistributionObjects, strings.TrimRight(strings.Repeat("?", len(ids)), ","))
	rows, err := m.dbm.Query(ctx, m.dbm.Rebind(ctx, q), utils.ArrayInt64(ids).ToArrayInterface()...)
	if err != nil {
		if err == sql.ErrNoRows {
			return o, nil
		}
		return nil, err
	}
	for rows.Next() {
		var obj Object
		err = rows.Scan(
			&obj.Id,
			&obj.DistributionId,
			&obj.RecipientId,
			&obj.RespondentId,
			&obj.PublishingStatus,
			&obj.PublishingLog,
			&obj.CreatedAt,
			&obj.UpdatedAt,
			&obj.PublishedAt,
		)
		if err != nil {
			return nil, err
		}
		o = append(o, &obj)
	}
	return o, nil
}

func (m *manager) FindObjectByLinkIds(ctx context.Context, ids ...int64) ([]*Object, error) {
	panic("implement me")
}

func (m *manager) FindObjectByRespondentAndLinkId(ctx context.Context, respondentId, id int64) ([]*Object, error) {
	panic("implement me")
}

func (m *manager) FindAllObjects(ctx context.Context, page, limit int) (items []*Object, total int, err error) {
	var (
		rows database.Rows
	)
	q := fmt.Sprintf(`SELECT count(id) FROM %s`, tableDistributionObjects)
	items = make([]*Object, 0)
	err = m.dbm.QueryRowAndBind(ctx, q, nil, &total)
	if err != nil || total < 1 {
		err = errors.Wrap(err, "It looks like the data is not exist")
		return
	}
	rows, err = m.dbm.Query(
		ctx, fmt.Sprintf(
			`SELECT 
						id, distribution_id, recipient_id, respondent_id,
						publishing_status, publishing_log, retry_count, 
						created_at, updated_at, published_at
					FROM %s LIMIT %d OFFSET %d`,
			tableDistributionObjects, limit, (page-1)*limit),
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return items, total, nil
		}
		return
	}
	for rows.Next() {
		var (
			item Object
		)
		err = rows.Scan(
			&item.Id,
			&item.DistributionId,
			&item.RecipientId,
			&item.RespondentId,
			&item.PublishingStatus,
			&item.PublishingLog,
			&item.RetryCount,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.PublishedAt,
		)
		if err != nil {
			return
		}
		items = append(items, &item)
	}
	return
}

func (m *manager) InsertObjects(ctx context.Context, items []*Object) error {
	q := fmt.Sprintf(`
		INSERT INTO %s (
			distribution_id, recipient_id, respondent_id, link_id, 
			publishing_status, publishing_log, created_by,
			retry_count, created_at
		) VALUES`, tableDistributionObjects)
	placeholders := make([]string, 0)
	values := make([]interface{}, 0)
	for _, item := range items {
		placeholders = append(
			placeholders,
			`(?, ?, ?, ?, 'none', '[]', ?, 0, NOW())`,
		)
		values = append(
			values, item.DistributionId, item.RecipientId, item.RespondentId,
			item.LinkId, item.CreatedBy,
		)
	}
	q = m.dbm.Rebind(ctx, fmt.Sprintf(`%s %s`, q, strings.Join(placeholders, ",")))
	cmd, err2 := m.dbm.Exec(ctx, q, values...)
	if err2 != nil {
		return errors.Wrap(err2, "failed saving groups. some errors in constraint or data.")
	}
	if cmd.RowsAffected() > 0 {
		return nil
	}
	return fmt.Errorf("no rows created")
}

func (m *manager) UpdateObject(ctx context.Context, item Object) error {
	if item.Id < 1 {
		return fmt.Errorf("please provide the correct identifier")
	}
	args := []interface{}{
		item.DistributionId,
		item.RecipientId,
		item.RespondentId,
		item.PublishingStatus,
		item.PublishingLog,
		item.RetryCount,
		item.UpdatedBy,
	}
	var publishedAt interface{} = item.PublishedAt
	if item.PublishingStatus != PublishingNone {
		publishedAt = "NOW()"
	}
	args = append(args, publishedAt)
	args = append(args, item.Id)
	q := fmt.Sprintf(`
		UPDATE %s 
		SET 
			distribution_id = ?,
			recipient_id = ?,
			respondent_id = ?,
			publishing_status = ?,
			publishing_log = ?,
			retry_count = ?,
			updated_by = ?,
			updated_at = NOW(),
			published_at = ?
		WHERE id = ?`, tableDistributionObjects)
	q = m.dbm.Rebind(ctx, q)
	cmd, err2 := m.dbm.Exec(ctx, q, args...)
	if err2 != nil {
		return err2
	}
	if cmd.RowsAffected() > 0 {
		return nil
	}
	return fmt.Errorf("no rows updated")
}

func (m *manager) DeleteObjectByIds(ctx context.Context, ids ...int64) error {
	q := m.dbm.Rebind(ctx, fmt.Sprintf(`
			DELETE FROM %s WHERE id IN (%s)
		`, tableDistributionObjects, strings.TrimRight(strings.Repeat("?,", len(ids)), ",")))
	rs, err := m.dbm.Exec(
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

