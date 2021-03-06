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
	tableDistribution          = "distributions"
	tableDistributionObjects   = "distribution_objects"
	tableDistributionMailQueue = "distribution_mail_queue"
	tableDistributionLog       = "distribution_log"
)

func NewDistributionDomain(dbm database.IManager) IManager {
	return &manager{dbm: dbm}
}

func (m *manager) FindByIds(ctx context.Context, ids ...int64) ([]*Distribution, error) {
	d := make([]*Distribution, 0)
	q := fmt.Sprintf(`
		SELECT 
			id, topic, created_by, range_start, range_end, disabled, archived, 
			distributed, distribution_limit, distribution_count,  
			created_at, updated_at, disabled_at, archived_at, distributed_at 
		FROM %s WHERE id IN (%s)`, tableDistribution, strings.TrimRight(strings.Repeat("?,", len(ids)), ","))
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
		disabledAt    interface{} = item.DisabledAt
		archivedAt    interface{} = item.ArchivedAt
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

func (m *manager) UpdateStatusAndCountByIds(ctx context.Context, ids ...int64) error {
	if len(ids) < 1 {
		return fmt.Errorf("please provide the correct identifier")
	}
	q := fmt.Sprintf(`
		UPDATE %s 
		SET 
			distributed = true,
			distribution_count = distribution_count + 1
		WHERE id IN (%s)`,
		tableDistribution, strings.TrimRight(strings.Repeat("?,", len(ids)), ","),
	)
	q = m.dbm.Rebind(ctx, q)
	cmd, err := m.dbm.Exec(ctx, q, utils.ArrayInt64(ids).ToArrayInterface()...)
	if err != nil {
		return err
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
	items := make([]*Object, 0)
	q := fmt.Sprintf(`
		SELECT 
			id, distribution_id, recipient_id, respondent_id, publishing_status, publishing_log, retry_count,
			created_by, updated_by, link_id, created_at, updated_at, published_at 
		FROM %s WHERE id IN (%s)`, tableDistributionObjects, strings.TrimRight(strings.Repeat("?,", len(ids)), ","))
	rows, err := m.dbm.Query(ctx, m.dbm.Rebind(ctx, q), utils.ArrayInt64(ids).ToArrayInterface()...)
	if err != nil {
		if err == sql.ErrNoRows {
			return items, nil
		}
		return nil, err
	}
	for rows.Next() {
		var (
			item              Object
			updatedBy, linkId, retryCount sql.NullInt64
		)
		err = rows.Scan(
			&item.Id,
			&item.DistributionId,
			&item.RecipientId,
			&item.RespondentId,
			&item.PublishingStatus,
			&item.PublishingLog,
			&retryCount,
			&item.CreatedBy,
			&updatedBy,
			&linkId,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.PublishedAt,
		)
		item.UpdatedBy = updatedBy.Int64
		item.LinkId = linkId.Int64
		item.RetryCount = retryCount.Int64
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

func (m *manager) FindObjectsByDistributionIds(ctx context.Context, ids ...int64) ([]*Object, error) {
	items := make([]*Object, 0)
	q := fmt.Sprintf(`
		SELECT 
			id, distribution_id, recipient_id, respondent_id, publishing_status, publishing_log, retry_count,
			created_by, updated_by, link_id, created_at, updated_at, published_at 
		FROM %s WHERE distribution_id IN (%s)`, tableDistributionObjects, strings.TrimRight(strings.Repeat("?,", len(ids)), ","))
	rows, err := m.dbm.Query(ctx, m.dbm.Rebind(ctx, q), utils.ArrayInt64(ids).ToArrayInterface()...)
	if err != nil {
		if err == sql.ErrNoRows {
			return items, nil
		}
		return nil, err
	}
	for rows.Next() {
		var (
			item              Object
			retryCount, updatedBy, linkId sql.NullInt64
		)
		err = rows.Scan(
			&item.Id,
			&item.DistributionId,
			&item.RecipientId,
			&item.RespondentId,
			&item.PublishingStatus,
			&item.PublishingLog,
			&retryCount,
			&item.CreatedBy,
			&updatedBy,
			&linkId,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.PublishedAt,
		)
		if err != nil {
			return nil, err
		}
		item.RetryCount = retryCount.Int64
		item.UpdatedBy = updatedBy.Int64
		item.LinkId = linkId.Int64
		items = append(items, &item)
	}
	return items, nil
}

func (m *manager) FindObjectByLinkIds(ctx context.Context, ids ...int64) ([]*Object, error) {
	items := make([]*Object, 0)
	q := fmt.Sprintf(`
		SELECT 
			id, distribution_id, recipient_id, respondent_id, publishing_status, publishing_log, retry_count,
			created_by, updated_by, link_id, created_at, updated_at, published_at 
		FROM %s WHERE link_id IN (%s)`, tableDistributionObjects, strings.TrimRight(strings.Repeat("?,", len(ids)), ","))
	rows, err := m.dbm.Query(ctx, m.dbm.Rebind(ctx, q), utils.ArrayInt64(ids).ToArrayInterface()...)
	if err != nil {
		if err == sql.ErrNoRows {
			return items, nil
		}
		return nil, err
	}
	for rows.Next() {
		var (
			item              Object
			updatedBy, linkId, retryCount sql.NullInt64
		)
		err = rows.Scan(
			&item.Id,
			&item.DistributionId,
			&item.RecipientId,
			&item.RespondentId,
			&item.PublishingStatus,
			&item.PublishingLog,
			&retryCount,
			&item.CreatedBy,
			&updatedBy,
			&linkId,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.PublishedAt,
		)
		item.UpdatedBy = updatedBy.Int64
		item.LinkId = linkId.Int64
		item.RetryCount = retryCount.Int64
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

func (m *manager) FindObjectByRespondentAndLinkId(ctx context.Context, respondentId, id int64) ([]*Object, error) {
	items := make([]*Object, 0)
	q := fmt.Sprintf(`
		SELECT 
			id, distribution_id, recipient_id, respondent_id, publishing_status, publishing_log, retry_count,
			created_by, updated_by, link_id, created_at, updated_at, published_at 
		FROM %s WHERE respondent_id = ? AND link_id = ?`, tableDistributionObjects)
	rows, err := m.dbm.Query(ctx, m.dbm.Rebind(ctx, q), respondentId, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return items, nil
		}
		return nil, err
	}
	for rows.Next() {
		var (
			item              Object
			updatedBy, linkId, retryCount sql.NullInt64
		)
		err = rows.Scan(
			&item.Id,
			&item.DistributionId,
			&item.RecipientId,
			&item.RespondentId,
			&item.PublishingStatus,
			&item.PublishingLog,
			&retryCount,
			&item.CreatedBy,
			&updatedBy,
			&linkId,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.PublishedAt,
		)
		item.UpdatedBy = updatedBy.Int64
		item.LinkId = linkId.Int64
		item.RetryCount = retryCount.Int64
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
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
						created_by, updated_by, link_id,
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
			updatedBy, linkId, retryCount sql.NullInt64
		)
		err = rows.Scan(
			&item.Id,
			&item.DistributionId,
			&item.RecipientId,
			&item.RespondentId,
			&item.PublishingStatus,
			&item.PublishingLog,
			&retryCount,
			&item.CreatedBy,
			&updatedBy,
			&linkId,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.PublishedAt,
		)
		item.UpdatedBy = updatedBy.Int64
		item.LinkId = linkId.Int64
		item.RetryCount = retryCount.Int64
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

func (m *manager) UpdateObjectStatusByIds(ctx context.Context, status PublishingStatus, ids ...int64) error {
	if len(ids) < 1 {
		return fmt.Errorf("please provide the correct identifier")
	}
	var publishedArg = ""
	if status != PublishingNone {
		publishedArg = ",published_at = NOW()"
	}
	q := fmt.Sprintf(`
		UPDATE %s 
		SET 
			publishing_status = ?,
			updated_at = NOW()
			%s
		WHERE id IN (%s)`,
		tableDistributionObjects, publishedArg, strings.TrimRight(strings.Repeat("?,", len(ids)), ","),
	)
	q = m.dbm.Rebind(ctx, q)
	cmd, err2 := m.dbm.Exec(ctx, q, status, utils.ArrayInt64(ids).ToArrayInterface())
	if err2 != nil {
		return err2
	}
	if cmd.RowsAffected() > 0 {
		return nil
	}
	return fmt.Errorf("no rows updated")
}

func (m *manager) UpdateObjectRetryCountByIds(ctx context.Context, ids ...int64) error {
	if len(ids) < 1 {
		return fmt.Errorf("please provide the correct identifier")
	}
	q := fmt.Sprintf(`
		UPDATE %s 
		SET 
			retry_count = retry_count + 1
		WHERE id IN (%s)`,
		tableDistributionObjects, strings.TrimRight(strings.Repeat("?,", len(ids)), ","),
	)
	q = m.dbm.Rebind(ctx, q)
	cmd, err2 := m.dbm.Exec(ctx, q, ids)
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

func (m *manager) FindAllQueues(ctx context.Context, page, limit int) (items []*Queue, total int64, err error) {
	var (
		rows database.Rows
	)
	q := fmt.Sprintf(`SELECT count(id) FROM %s`, tableDistributionMailQueue)
	items = make([]*Queue, 0)
	err = m.dbm.QueryRowAndBind(ctx, q, nil, &total)
	if err != nil || total < 1 {
		err = errors.Wrap(err, "It looks like the data is not exist")
		return
	}
	rows, err = m.dbm.Query(
		ctx, fmt.Sprintf(
			`SELECT 
						id, distribution_object_id, recipient_id, respondent_id, 
						from_email, to_email, subject, template, arguments
					FROM %s LIMIT %d OFFSET %d`,
			tableDistributionMailQueue, limit, (page-1)*limit),
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return items, total, nil
		}
		return
	}
	for rows.Next() {
		var (
			item Queue
		)
		err = rows.Scan(
			&item.Id,
			&item.DistributionObjectId,
			&item.RecipientId,
			&item.RespondentId,
			&item.FromEmail,
			&item.ToEmail,
			&item.Subject,
			&item.Template,
			&item.Arguments,
		)
		if err != nil {
			return
		}
		items = append(items, &item)
	}
	return
}

func (m *manager) InsertQueues(ctx context.Context, items []*Queue) ([]int64, error) {
	q := fmt.Sprintf(`
		INSERT INTO %s (
			distribution_object_id, recipient_id, respondent_id, 
			from_email, to_email, subject, template, arguments
		) VALUES`, tableDistributionMailQueue)
	placeholders := make([]string, 0)
	values := make([]interface{}, 0)
	ids := make([]int64, 0)
	for _, item := range items {
		placeholders = append(
			placeholders,
			`(?, ?, ?, ?, ?, ?, ?, ?)`,
		)
		values = append(
			values, item.DistributionObjectId, item.RecipientId, item.RespondentId,
			item.FromEmail, item.ToEmail, item.Subject, item.Template, item.Arguments,
		)
	}
	q = m.dbm.Rebind(ctx, fmt.Sprintf(`%s %s RETURNING id`, q, strings.Join(placeholders, ",")))
	rows, err := m.dbm.Query(ctx, q, values...)
	if err != nil {
		return ids, errors.Wrap(err, "failed saving queue. some errors in constraint or data.")
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

func (m *manager) DeleteQueueByIds(ctx context.Context, ids ...int64) error {
	q := m.dbm.Rebind(ctx, fmt.Sprintf(`
			DELETE FROM %s WHERE id IN (%s)
		`, tableDistributionMailQueue, strings.TrimRight(strings.Repeat("?,", len(ids)), ",")))
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

func (m *manager) InsertLogs(ctx context.Context, items []*Log) error {
	q := fmt.Sprintf(`
		INSERT INTO %s (
			action, values, values_prev, 
			notes, at
		) VALUES`, tableDistributionLog)
	placeholders := make([]string, 0)
	values := make([]interface{}, 0)
	for _, item := range items {
		placeholders = append(
			placeholders,
			`(?, ?, ?, ?, NOW())`,
		)
		values = append(
			values, item.Action, item.Values, item.ValuesPrev, item.Notes,
		)
	}
	q = m.dbm.Rebind(ctx, fmt.Sprintf(`%s %s`, q, strings.Join(placeholders, ",")))
	cmd, err := m.dbm.Exec(ctx, q, values...)
	if err != nil {
		return errors.Wrap(err, "failed saving logs. some errors in constraint or data.")
	}
	if cmd.RowsAffected() < 1 {
		return fmt.Errorf("no rows created")
	}
	return nil
}
