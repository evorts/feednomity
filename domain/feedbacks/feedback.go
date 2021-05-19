package feedbacks

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/utils"
	"github.com/pkg/errors"
	"strings"
)

const (
	tableFeedback    = "feedbacks"
	tableFeedbackLog = "feedback_log"
)

type manager struct {
	dbm database.IManager
}

type IFeedback interface {
	FindByIds(ctx context.Context, ids ...int64) ([]*Feedback, error)
	FindByDistId(ctx context.Context, distId, distObjId int64) ([]*Feedback, error)
	FindItem(ctx context.Context, distId, distObjId, recipientId, respondentId int64) (*Feedback, error)
	FindAll(ctx context.Context, page, limit int) (items []*Feedback, total int, err error)

	InsertMultiple(ctx context.Context, items []*Feedback) error
	Update(ctx context.Context, item Feedback) error
	UpdateStatusAndContent(ctx context.Context, id int64, status Status, content map[string]interface{}) error
	DeleteByIds(ctx context.Context, ids ...int64) error
}

func NewFeedbackDomain(dbm database.IManager) IFeedback {
	return &manager{dbm: dbm}
}

func (m *manager) FindByIds(ctx context.Context, ids ...int64) ([]*Feedback, error) {
	items := make([]*Feedback, 0)
	q := fmt.Sprintf(`
		SELECT 
			id, distribution_id, distribution_topic, distribution_object_id, range_start, range_end,
			respondent_id, respondent_username, respondent_name, respondent_email, 
			respondent_group_id, respondent_group_name, respondent_org_id, respondent_org_name,
			recipient_id, recipient_username, recipient_name, recipient_email, 
			recipient_group_id, recipient_group_name, recipient_org_id, recipient_org_name,
			link_id, hash, status, content, created_at, updated_at 
		FROM %s WHERE id IN (%s)`, tableFeedback, strings.TrimRight(strings.Repeat("?,", len(ids)), ","))
	rows, err := m.dbm.Query(ctx, m.dbm.Rebind(ctx, q), utils.ArrayInt64(ids).ToArrayInterface()...)
	if err != nil {
		if err == sql.ErrNoRows {
			return items, nil
		}
		return nil, err
	}
	for rows.Next() {
		var (
			item Feedback
		)
		err = rows.Scan(
			&item.Id,
			&item.DistributionId,
			&item.DistributionTopic,
			&item.DistributionObjectId,
			&item.RangeStart,
			&item.RangeEnd,
			&item.RespondentId,
			&item.RespondentUsername,
			&item.RespondentName,
			&item.RespondentEmail,
			&item.RespondentGroupId,
			&item.RespondentGroupName,
			&item.RespondentOrgId,
			&item.RespondentOrgName,
			&item.RecipientId,
			&item.RecipientUsername,
			&item.RecipientName,
			&item.RecipientEmail,
			&item.RecipientGroupId,
			&item.RecipientGroupName,
			&item.RecipientOrgId,
			&item.RecipientOrgName,
			&item.LinkId,
			&item.Hash,
			&item.Status,
			&item.Content,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

func (m *manager) FindByDistId(ctx context.Context, distId, distObjId int64) ([]*Feedback, error) {
	items := make([]*Feedback, 0)
	q := fmt.Sprintf(`
		SELECT 
			id, distribution_id, distribution_topic, distribution_object_id, range_start, range_end,
			respondent_id, respondent_username, respondent_name, respondent_email, 
			respondent_group_id, respondent_group_name, respondent_org_id, respondent_org_name,
			recipient_id, recipient_username, recipient_name, recipient_email, 
			recipient_group_id, recipient_group_name, recipient_org_id, recipient_org_name,
			link_id, hash, status, content, created_at, updated_at 
		FROM %s WHERE distribution_id = ? AND distribution_object_id = ?`, tableFeedback,
	)
	rows, err := m.dbm.Query(
		ctx,
		m.dbm.Rebind(ctx, q),
		[]interface{}{distId, distObjId}...,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return items, nil
		}
		return nil, err
	}
	for rows.Next() {
		var (
			item Feedback
		)
		err = rows.Scan(
			&item.Id,
			&item.DistributionId,
			&item.DistributionTopic,
			&item.DistributionObjectId,
			&item.RangeStart,
			&item.RangeEnd,
			&item.RespondentId,
			&item.RespondentUsername,
			&item.RespondentName,
			&item.RespondentEmail,
			&item.RespondentGroupId,
			&item.RespondentGroupName,
			&item.RespondentOrgId,
			&item.RespondentOrgName,
			&item.RecipientId,
			&item.RecipientUsername,
			&item.RecipientName,
			&item.RecipientEmail,
			&item.RecipientGroupId,
			&item.RecipientGroupName,
			&item.RecipientOrgId,
			&item.RecipientOrgName,
			&item.LinkId,
			&item.Hash,
			&item.Status,
			&item.Content,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

func (m *manager) FindItem(ctx context.Context, distId, distObjId, recipientId, respondentId int64) (*Feedback, error) {
	q := fmt.Sprintf(`
		SELECT 
			id, distribution_id, distribution_topic, distribution_object_id, range_start, range_end,
			respondent_id, respondent_username, respondent_name, respondent_email, 
			respondent_group_id, respondent_group_name, respondent_org_id, respondent_org_name,
			recipient_id, recipient_username, recipient_name, recipient_email, 
			recipient_group_id, recipient_group_name, recipient_org_id, recipient_org_name,
			link_id, hash, status, content, created_at, updated_at 
		FROM %s WHERE 
			distribution_id = ? 
			AND distribution_object_id = ?
			AND recipient_id = ?
			AND respondent_id = ?`, tableFeedback,
	)
	var item Feedback
	err := m.dbm.QueryRow(
		ctx,
		m.dbm.Rebind(ctx, q),
		[]interface{}{distId, distObjId, recipientId, respondentId}...,
	).Scan(
		&item.Id,
		&item.DistributionId,
		&item.DistributionTopic,
		&item.DistributionObjectId,
		&item.RangeStart,
		&item.RangeEnd,
		&item.RespondentId,
		&item.RespondentUsername,
		&item.RespondentName,
		&item.RespondentEmail,
		&item.RespondentGroupId,
		&item.RespondentGroupName,
		&item.RespondentOrgId,
		&item.RespondentOrgName,
		&item.RecipientId,
		&item.RecipientUsername,
		&item.RecipientName,
		&item.RecipientEmail,
		&item.RecipientGroupId,
		&item.RecipientGroupName,
		&item.RecipientOrgId,
		&item.RecipientOrgName,
		&item.LinkId,
		&item.Hash,
		&item.Status,
		&item.Content,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	return &item, err
}

func (m *manager) FindAll(ctx context.Context, page, limit int) (items []*Feedback, total int, err error) {
	var (
		rows database.Rows
	)
	q := fmt.Sprintf(`SELECT count(id) FROM %s`, tableFeedback)
	items = make([]*Feedback, 0)
	err = m.dbm.QueryRowAndBind(ctx, q, nil, &total)
	if err != nil || total < 1 {
		err = errors.Wrap(err, "It looks like the data is not exist")
		return
	}
	rows, err = m.dbm.Query(
		ctx, fmt.Sprintf(
			`SELECT 
						id, distribution_id, distribution_topic, distribution_object_id, range_start, range_end,
						respondent_id, respondent_username, respondent_name, respondent_email, 
						respondent_group_id, respondent_group_name, respondent_org_id, respondent_org_name,
						recipient_id, recipient_username, recipient_name, recipient_email, 
						recipient_group_id, recipient_group_name, recipient_org_id, recipient_org_name,
						link_id, hash, status, content, created_at, updated_at
					FROM %s LIMIT %d OFFSET %d`,
			tableFeedback, limit, (page-1)*limit),
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return items, total, nil
		}
		return
	}
	for rows.Next() {
		var (
			item Feedback
		)
		err = rows.Scan(
			&item.Id,
			&item.DistributionId,
			&item.DistributionTopic,
			&item.DistributionObjectId,
			&item.RangeStart,
			&item.RangeEnd,
			&item.RespondentId,
			&item.RespondentUsername,
			&item.RespondentName,
			&item.RespondentEmail,
			&item.RespondentGroupId,
			&item.RespondentGroupName,
			&item.RespondentOrgId,
			&item.RespondentOrgName,
			&item.RecipientId,
			&item.RecipientUsername,
			&item.RecipientName,
			&item.RecipientEmail,
			&item.RecipientGroupId,
			&item.RecipientGroupName,
			&item.RecipientOrgId,
			&item.RecipientOrgName,
			&item.LinkId,
			&item.Hash,
			&item.Status,
			&item.Content,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return
		}
		items = append(items, &item)
	}
	return
}

func (m *manager) InsertMultiple(ctx context.Context, items []*Feedback) error {
	q := fmt.Sprintf(`
		INSERT INTO %s (
			distribution_id, distribution_topic, distribution_object_id, range_start, range_end,
			respondent_id, respondent_username, respondent_name, respondent_email, 
			respondent_group_id, respondent_group_name, respondent_org_id, respondent_org_name,
			recipient_id, recipient_username, recipient_name, recipient_email, 
			recipient_group_id, recipient_group_name, recipient_org_id, recipient_org_name,
			link_id, hash, status, content, created_at
		) VALUES`, tableFeedback)
	placeholders := make([]string, 0)
	values := make([]interface{}, 0)
	for _, item := range items {
		placeholders = append(
			placeholders,
			`(
				?, ?, ?, ?, ?, 
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?, NOW()
			)`,
		)
		values = append(values,
			&item.DistributionId, &item.DistributionTopic, &item.DistributionObjectId, &item.RangeStart, &item.RangeEnd,
			&item.RespondentId, &item.RespondentUsername, &item.RespondentName, &item.RespondentEmail,
			&item.RespondentGroupId, &item.RespondentGroupName, &item.RespondentOrgId, &item.RespondentOrgName,
			&item.RecipientId, &item.RecipientUsername, &item.RecipientName, &item.RecipientEmail,
			&item.RecipientGroupId, &item.RecipientGroupName, &item.RecipientOrgId, &item.RecipientOrgName,
			&item.LinkId, &item.Hash, &item.Status, &item.Content,
		)
	}
	q = m.dbm.Rebind(ctx, fmt.Sprintf(`%s %s`, q, strings.Join(placeholders, ",")))
	cmd, err2 := m.dbm.Exec(ctx, q, values...)
	if err2 != nil {
		return errors.Wrap(err2, "failed saving feedbacks. some errors in constraint or data.")
	}
	if cmd.RowsAffected() > 0 {
		return nil
	}
	return fmt.Errorf("no rows created")
}

func (m *manager) Update(ctx context.Context, item Feedback) error {
	if item.Id < 1 {
		return fmt.Errorf("please provide the correct identifier")
	}
	args := []interface{}{
		&item.DistributionId, &item.DistributionTopic, &item.DistributionObjectId, &item.RangeStart, &item.RangeEnd,
		&item.RespondentId, &item.RespondentUsername, &item.RespondentName, &item.RespondentEmail,
		&item.RespondentGroupId, &item.RespondentGroupName, &item.RespondentOrgId, &item.RespondentOrgName,
		&item.RecipientId, &item.RecipientUsername, &item.RecipientName, &item.RecipientEmail,
		&item.RecipientGroupId, &item.RecipientGroupName, &item.RecipientOrgId, &item.RecipientOrgName,
		&item.LinkId, &item.Hash, &item.Status, &item.Content,
	}
	args = append(args, item.Id)
	q := fmt.Sprintf(`
		UPDATE %s 
		SET 
			distribution_id = ?, 
			distribution_topic = ?, 
			distribution_object_id = ?, 
			range_start = ?, 
			range_end = ?,
			respondent_id = ?, 
			respondent_username = ?, 
			respondent_name = ?, 
			respondent_email = ?, 
			respondent_group_id = ?, 
			respondent_group_name = ?, 
			respondent_org_id = ?, 
			respondent_org_name = ?,
			recipient_id = ?, 
			recipient_username = ?, 
			recipient_name = ?, 
			recipient_email = ?, 
			recipient_group_id = ?, 
			recipient_group_name = ?, 
			recipient_org_id = ?, 
			recipient_org_name = ?,
			link_id = ?, 
			hash = ?, 
			status = ?, 
			content = ?, 
			updated_at = NOW()
		WHERE id = ?`, tableFeedback)
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

func (m *manager) UpdateStatusAndContent(ctx context.Context, id int64, status Status, content map[string]interface{}) error {
	if id < 1 {
		return fmt.Errorf("please provide the correct identifier")
	}
	q := fmt.Sprintf(`
		UPDATE %s 
		SET 
			status = ?, 
			content = ?, 
			updated_at = NOW()
		WHERE id = ?`, tableFeedback)
	q = m.dbm.Rebind(ctx, q)
	cmd, err2 := m.dbm.Exec(ctx, q, status, content, id)
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
		`, tableFeedback, strings.TrimRight(strings.Repeat("?,", len(ids)), ",")))
	rs, err := m.dbm.Exec(ctx, q, utils.ArrayInt64(ids).ToArrayInterface()...)
	if err != nil {
		return err
	}
	if rs.RowsAffected() < 1 {
		return errors.New("not a single record removed")
	}
	return nil
}
