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
	FindByDistId(ctx context.Context, distId int64) ([]*Feedback, error)
	FindByDistAndObjectId(ctx context.Context, distId, distObjId int64) ([]*Feedback, error)
	FindByGroupId(ctx context.Context, id int64, page, limit int) (items []*Feedback, total int, err error)
	FindByOrgId(ctx context.Context, id int64, page, limit int) (items []*Feedback, total int, err error)
	FindByRespondentId(ctx context.Context, id int64, page, limit int) (items []*Feedback, total int, err error)
	FindByRecipientId(ctx context.Context, id int64, page, limit int) (items []*Feedback, total int, err error)
	FindAllWithFilter(ctx context.Context, filter map[string]interface{}, withContent bool) (items []*Feedback, total int, err error)
	FindItem(ctx context.Context, distId, distObjId, recipientId, respondentId int64) (*Feedback, error)
	FindAll(ctx context.Context, page, limit int) (items []*Feedback, total int, err error)
	SummaryByDistribution(ctx context.Context, page, limit int, filter map[string]interface{}) (items []*Feedback, total int, err error)

	InsertMultiple(ctx context.Context, items []*Feedback) error
	// UpsertMultiple success items format [ [feedbackId, distributionId, distObjectId, respondentId, recipientId] ]
	UpsertMultiple(ctx context.Context, items []*Feedback) (successItems [][]int64, err error)
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
			respondent_role, respondent_assignment,
			recipient_id, recipient_username, recipient_name, recipient_email, 
			recipient_group_id, recipient_group_name, recipient_org_id, recipient_org_name,
			recipient_role, recipient_assignment,
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
			&item.Id, &item.DistributionId, &item.DistributionTopic, &item.DistributionObjectId, &item.RangeStart, &item.RangeEnd,
			&item.RespondentId, &item.RespondentUsername, &item.RespondentName, &item.RespondentEmail,
			&item.RespondentGroupId, &item.RespondentGroupName, &item.RespondentOrgId, &item.RespondentOrgName,
			&item.RespondentRole, &item.RespondentAssignment,
			&item.RecipientId, &item.RecipientUsername, &item.RecipientName, &item.RecipientEmail,
			&item.RecipientGroupId, &item.RecipientGroupName, &item.RecipientOrgId, &item.RecipientOrgName,
			&item.RecipientRole, &item.RecipientAssignment,
			&item.LinkId, &item.Hash, &item.Status, &item.Content, &item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

func (m *manager) FindByDistId(ctx context.Context, distId int64) ([]*Feedback, error) {
	items := make([]*Feedback, 0)
	q := fmt.Sprintf(`
		SELECT 
			id, distribution_id, distribution_topic, distribution_object_id, range_start, range_end,
			respondent_id, respondent_username, respondent_name, respondent_email, 
			respondent_group_id, respondent_group_name, respondent_org_id, respondent_org_name,
			respondent_role, respondent_assignment,
			recipient_id, recipient_username, recipient_name, recipient_email, 
			recipient_group_id, recipient_group_name, recipient_org_id, recipient_org_name,
			recipient_role, recipient_assignment,
			link_id, hash, status, content, created_at, updated_at 
		FROM %s WHERE distribution_id = ?`, tableFeedback,
	)
	rows, err := m.dbm.Query(
		ctx,
		m.dbm.Rebind(ctx, q),
		[]interface{}{distId}...,
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
			&item.RespondentRole,
			&item.RespondentAssignment,
			&item.RecipientId,
			&item.RecipientUsername,
			&item.RecipientName,
			&item.RecipientEmail,
			&item.RecipientGroupId,
			&item.RecipientGroupName,
			&item.RecipientOrgId,
			&item.RecipientOrgName,
			&item.RecipientRole,
			&item.RecipientAssignment,
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

func (m *manager) FindByDistAndObjectId(ctx context.Context, distId, distObjId int64) ([]*Feedback, error) {
	items := make([]*Feedback, 0)
	q := fmt.Sprintf(`
		SELECT 
			id, distribution_id, distribution_topic, distribution_object_id, range_start, range_end,
			respondent_id, respondent_username, respondent_name, respondent_email, 
			respondent_group_id, respondent_group_name, respondent_org_id, respondent_org_name,
			respondent_role, respondent_assignment,
			recipient_id, recipient_username, recipient_name, recipient_email, 
			recipient_group_id, recipient_group_name, recipient_org_id, recipient_org_name,
			recipient_role, recipient_assignment,
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
			&item.RespondentRole,
			&item.RespondentAssignment,
			&item.RecipientId,
			&item.RecipientUsername,
			&item.RecipientName,
			&item.RecipientEmail,
			&item.RecipientGroupId,
			&item.RecipientGroupName,
			&item.RecipientOrgId,
			&item.RecipientOrgName,
			&item.RecipientRole,
			&item.RecipientAssignment,
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

func (m *manager) FindByGroupId(ctx context.Context, id int64, page, limit int) (items []*Feedback, total int, err error) {
	var rows database.Rows
	items = make([]*Feedback, 0)
	q := fmt.Sprintf(`SELECT count(id) FROM %s WHERE recipient_group_id = $1`, tableFeedback)
	err = m.dbm.QueryRowAndBind(ctx, q, []interface{}{id}, &total)
	if err != nil || total < 1 {
		err = errors.Wrap(err, "It looks like the data is not exist")
		return
	}
	items = make([]*Feedback, 0)
	q = fmt.Sprintf(`
		SELECT 
			id, distribution_id, distribution_topic, distribution_object_id, range_start, range_end,
			respondent_id, respondent_username, respondent_name, respondent_email, 
			respondent_group_id, respondent_group_name, respondent_org_id, respondent_org_name,
			respondent_role, respondent_assignment,
			recipient_id, recipient_username, recipient_name, recipient_email, 
			recipient_group_id, recipient_group_name, recipient_org_id, recipient_org_name,
			recipient_role, recipient_assignment,
			link_id, hash, status, content, created_at, updated_at 
		FROM %s 
		JOIN (VALUES ('%s'::feedback_status, 1), ('%s'::feedback_status, 2), ('%s'::feedback_status, 3)) AS x(value, order_number) on %s.status = x.value
		WHERE recipient_group_id = ? 
		ORDER BY distribution_id desc, x.order_number, recipient_name asc, created_at, updated_at desc
		LIMIT %d OFFSET %d`, tableFeedback, StatusDraft, StatusNotStarted, StatusFinal, tableFeedback, limit, (page-1)*limit,
	)
	rows, err = m.dbm.Query(ctx, m.dbm.Rebind(ctx, q), id)
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
			&item.RespondentRole,
			&item.RespondentAssignment,
			&item.RecipientId,
			&item.RecipientUsername,
			&item.RecipientName,
			&item.RecipientEmail,
			&item.RecipientGroupId,
			&item.RecipientGroupName,
			&item.RecipientOrgId,
			&item.RecipientOrgName,
			&item.RecipientRole,
			&item.RecipientAssignment,
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

func (m *manager) FindByOrgId(ctx context.Context, id int64, page, limit int) (items []*Feedback, total int, err error) {
	var rows database.Rows
	items = make([]*Feedback, 0)
	q := fmt.Sprintf(`SELECT count(id) FROM %s WHERE recipient_org_id = $1`, tableFeedback)
	err = m.dbm.QueryRowAndBind(ctx, q, []interface{}{id}, &total)
	if err != nil || total < 1 {
		err = errors.Wrap(err, "It looks like the data is not exist")
		return
	}
	items = make([]*Feedback, 0)
	q = fmt.Sprintf(`
		SELECT 
			id, distribution_id, distribution_topic, distribution_object_id, range_start, range_end,
			respondent_id, respondent_username, respondent_name, respondent_email, 
			respondent_group_id, respondent_group_name, respondent_org_id, respondent_org_name,
			respondent_role, respondent_assignment,
			recipient_id, recipient_username, recipient_name, recipient_email, 
			recipient_group_id, recipient_group_name, recipient_org_id, recipient_org_name,
			recipient_role, recipient_assignment,
			link_id, hash, status, content, created_at, updated_at 
		FROM %s
		JOIN (VALUES ('%s'::feedback_status, 1), ('%s'::feedback_status, 2), ('%s'::feedback_status, 3)) AS x(value, order_number) on %s.status = x.value
		WHERE recipient_org_id = ? 
		ORDER BY distribution_id desc, x.order_number, recipient_name asc, created_at, updated_at desc
		LIMIT %d OFFSET %d`, tableFeedback, StatusDraft, StatusNotStarted, StatusFinal, tableFeedback, limit, (page-1)*limit,
	)
	rows, err = m.dbm.Query(ctx, m.dbm.Rebind(ctx, q), id)
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
			&item.RespondentRole,
			&item.RespondentAssignment,
			&item.RecipientId,
			&item.RecipientUsername,
			&item.RecipientName,
			&item.RecipientEmail,
			&item.RecipientGroupId,
			&item.RecipientGroupName,
			&item.RecipientOrgId,
			&item.RecipientOrgName,
			&item.RecipientRole,
			&item.RecipientAssignment,
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

func (m *manager) FindByRespondentId(ctx context.Context, id int64, page, limit int) (items []*Feedback, total int, err error) {
	var rows database.Rows
	items = make([]*Feedback, 0)
	q := fmt.Sprintf(`SELECT count(id) FROM %s WHERE respondent_id = $1`, tableFeedback)
	err = m.dbm.QueryRowAndBind(ctx, q, []interface{}{id}, &total)
	if err != nil || total < 1 {
		err = errors.Wrap(err, "It looks like the data is not exist")
		return
	}
	items = make([]*Feedback, 0)
	q = fmt.Sprintf(`
		SELECT 
			id, distribution_id, distribution_topic, distribution_object_id, range_start, range_end,
			respondent_id, respondent_username, respondent_name, respondent_email, 
			respondent_group_id, respondent_group_name, respondent_org_id, respondent_org_name,
			respondent_role, respondent_assignment,
			recipient_id, recipient_username, recipient_name, recipient_email, 
			recipient_group_id, recipient_group_name, recipient_org_id, recipient_org_name,
			recipient_role, recipient_assignment,
			link_id, hash, status, content, created_at, updated_at 
		FROM %s
		JOIN (VALUES ('%s'::feedback_status, 1), ('%s'::feedback_status, 2), ('%s'::feedback_status, 3)) AS x(value, order_number) on %s.status = x.value
		WHERE respondent_id = ?
		ORDER BY distribution_id desc, x.order_number, recipient_name asc, created_at, updated_at desc 
		LIMIT %d OFFSET %d
		`, tableFeedback, StatusDraft, StatusNotStarted, StatusFinal, tableFeedback, limit, (page-1)*limit,
	)
	rows, err = m.dbm.Query(ctx, m.dbm.Rebind(ctx, q), id)
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
			&item.Id, &item.DistributionId, &item.DistributionTopic, &item.DistributionObjectId, &item.RangeStart, &item.RangeEnd,
			&item.RespondentId, &item.RespondentUsername, &item.RespondentName, &item.RespondentEmail,
			&item.RespondentGroupId, &item.RespondentGroupName, &item.RespondentOrgId, &item.RespondentOrgName,
			&item.RespondentRole, &item.RespondentAssignment,
			&item.RecipientId, &item.RecipientUsername, &item.RecipientName, &item.RecipientEmail,
			&item.RecipientGroupId, &item.RecipientGroupName, &item.RecipientOrgId, &item.RecipientOrgName,
			&item.RecipientRole, &item.RecipientAssignment,
			&item.LinkId, &item.Hash, &item.Status, &item.Content, &item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return
		}
		items = append(items, &item)
	}
	return
}

func (m *manager) FindByRecipientId(ctx context.Context, id int64, page, limit int) (items []*Feedback, total int, err error) {
	var rows database.Rows
	items = make([]*Feedback, 0)
	q := fmt.Sprintf(`SELECT count(id) FROM %s WHERE recipient_id = $1`, tableFeedback)
	err = m.dbm.QueryRowAndBind(ctx, q, []interface{}{id}, &total)
	if err != nil || total < 1 {
		err = errors.Wrap(err, "It looks like the data is not exist")
		return
	}
	items = make([]*Feedback, 0)
	q = fmt.Sprintf(`
		SELECT 
			id, distribution_id, distribution_topic, distribution_object_id, range_start, range_end,
			respondent_id, respondent_username, respondent_name, respondent_email, 
			respondent_group_id, respondent_group_name, respondent_org_id, respondent_org_name,
			respondent_role, respondent_assignment,
			recipient_id, recipient_username, recipient_name, recipient_email, 
			recipient_group_id, recipient_group_name, recipient_org_id, recipient_org_name,
			recipient_role, recipient_assignment,
			link_id, hash, status, content, created_at, updated_at 
		FROM %s
		JOIN (VALUES ('%s'::feedback_status, 1), ('%s'::feedback_status, 2), ('%s'::feedback_status, 3)) AS x(value, order_number) on %s.status = x.value
		WHERE recipient_id = ?
		ORDER BY distribution_id desc, x.order_number, recipient_name asc, created_at, updated_at desc 
		LIMIT %d OFFSET %d
		`, tableFeedback, StatusDraft, StatusNotStarted, StatusFinal, tableFeedback, limit, (page-1)*limit,
	)
	rows, err = m.dbm.Query(ctx, m.dbm.Rebind(ctx, q), id)
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
			&item.Id, &item.DistributionId, &item.DistributionTopic, &item.DistributionObjectId, &item.RangeStart, &item.RangeEnd,
			&item.RespondentId, &item.RespondentUsername, &item.RespondentName, &item.RespondentEmail,
			&item.RespondentGroupId, &item.RespondentGroupName, &item.RespondentOrgId, &item.RespondentOrgName,
			&item.RespondentRole, &item.RespondentAssignment,
			&item.RecipientId, &item.RecipientUsername, &item.RecipientName, &item.RecipientEmail,
			&item.RecipientGroupId, &item.RecipientGroupName, &item.RecipientOrgId, &item.RecipientOrgName,
			&item.RecipientRole, &item.RecipientAssignment,
			&item.LinkId, &item.Hash, &item.Status, &item.Content, &item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return
		}
		items = append(items, &item)
	}
	return
}

func (m *manager) FindAllWithFilter(ctx context.Context, filter map[string]interface{}, withContent bool) (items []*Feedback, total int, err error) {
	var (
		where, qq = "", ""
		args      []interface{}
	)
	if filter != nil && len(filter) > 0 {
		qq, args = utils.GenerateFilters(filter)
		where = fmt.Sprintf(" WHERE %s", qq)
	}
	var rows database.Rows
	items = make([]*Feedback, 0)
	q := m.dbm.Rebind(ctx, fmt.Sprintf(`SELECT count(id) FROM %s %s`, tableFeedback, where))
	err = m.dbm.QueryRow(ctx, q, args...).Scan(&total)
	if err != nil || total < 1 {
		err = errors.Wrap(err, "It looks like the data is not exist")
		return
	}
	items = make([]*Feedback, 0)
	q = fmt.Sprintf(`
		SELECT 
			id, distribution_id, distribution_topic, distribution_object_id, range_start, range_end,
			respondent_id, respondent_username, respondent_name, respondent_email, 
			respondent_group_id, respondent_group_name, respondent_org_id, respondent_org_name,
			respondent_role, respondent_assignment,
			recipient_id, recipient_username, recipient_name, recipient_email, 
			recipient_group_id, recipient_group_name, recipient_org_id, recipient_org_name,
			recipient_role, recipient_assignment, 
			link_id, hash, status, content, created_at, updated_at %s
		FROM %s
		JOIN (VALUES ('%s'::feedback_status, 1), ('%s'::feedback_status, 2), ('%s'::feedback_status, 3)) AS x(value, order_number) on %s.status = x.value
		%s
		ORDER BY distribution_id desc, x.order_number, recipient_name asc, created_at, updated_at desc 
		`, utils.IIf(withContent, ",content", ""), tableFeedback, StatusDraft, StatusNotStarted, StatusFinal, tableFeedback, where,
	)
	rows, err = m.dbm.Query(ctx, m.dbm.Rebind(ctx, q), args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return items, total, nil
		}
		return
	}
	for rows.Next() {
		var (
			item Feedback
			dest []interface{}
		)
		dest = []interface{}{
			&item.Id, &item.DistributionId, &item.DistributionTopic, &item.DistributionObjectId, &item.RangeStart, &item.RangeEnd,
			&item.RespondentId, &item.RespondentUsername, &item.RespondentName, &item.RespondentEmail,
			&item.RespondentGroupId, &item.RespondentGroupName, &item.RespondentOrgId, &item.RespondentOrgName,
			&item.RespondentRole, &item.RespondentAssignment,
			&item.RecipientId, &item.RecipientUsername, &item.RecipientName, &item.RecipientEmail,
			&item.RecipientGroupId, &item.RecipientGroupName, &item.RecipientOrgId, &item.RecipientOrgName,
			&item.RecipientRole, &item.RecipientAssignment,
			&item.LinkId, &item.Hash, &item.Status, &item.Content, &item.CreatedAt, &item.UpdatedAt,
		}
		if withContent {
			dest = append(dest, &item.Content)
		}
		err = rows.Scan(dest...)
		if err != nil {
			return
		}
		items = append(items, &item)
	}
	return
}

func (m *manager) FindItem(ctx context.Context, distId, distObjId, recipientId, respondentId int64) (*Feedback, error) {
	q := fmt.Sprintf(`
		SELECT 
			id, distribution_id, distribution_topic, distribution_object_id, range_start, range_end,
			respondent_id, respondent_username, respondent_name, respondent_email, 
			respondent_group_id, respondent_group_name, respondent_org_id, respondent_org_name,
			respondent_role, respondent_assignment,
			recipient_id, recipient_username, recipient_name, recipient_email, 
			recipient_group_id, recipient_group_name, recipient_org_id, recipient_org_name,
			recipient_role, recipient_assignment,
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
		&item.RespondentRole,
		&item.RespondentAssignment,
		&item.RecipientId,
		&item.RecipientUsername,
		&item.RecipientName,
		&item.RecipientEmail,
		&item.RecipientGroupId,
		&item.RecipientGroupName,
		&item.RecipientOrgId,
		&item.RecipientOrgName,
		&item.RecipientRole,
		&item.RecipientAssignment,
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
						respondent_role, respondent_assignment,
						recipient_id, recipient_username, recipient_name, recipient_email, 
						recipient_group_id, recipient_group_name, recipient_org_id, recipient_org_name,
						recipient_role, recipient_assignment,
						link_id, hash, status, content, created_at, updated_at
					FROM %s
					JOIN (VALUES ('%s'::feedback_status, 1), ('%s'::feedback_status, 2), ('%s'::feedback_status, 3)) AS x(value, order_number) on %s.status = x.value
					ORDER BY distribution_id desc, x.order_number, recipient_name asc, created_at, updated_at desc
					LIMIT %d OFFSET %d
			`,
			tableFeedback, StatusDraft, StatusNotStarted, StatusFinal, tableFeedback, limit, (page-1)*limit),
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
			&item.RespondentRole,
			&item.RespondentAssignment,
			&item.RecipientId,
			&item.RecipientUsername,
			&item.RecipientName,
			&item.RecipientEmail,
			&item.RecipientGroupId,
			&item.RecipientGroupName,
			&item.RecipientOrgId,
			&item.RecipientOrgName,
			&item.RecipientRole,
			&item.RecipientAssignment,
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

func (m *manager) SummaryByDistribution(ctx context.Context, page, limit int, filter map[string]interface{}) (items []*Feedback, total int, err error) {
	var (
		rows      database.Rows
		where, qq = "", ""
		args      = make([]interface{}, 0)
	)
	q := fmt.Sprintf(`SELECT COUNT(DISTINCT distribution_id) as subtotal FROM %s GROUP BY distribution_id`, tableFeedback)
	err = m.dbm.QueryRowAndBind(ctx, q, nil, &total)
	if err != nil || total < 1 {
		err = errors.Wrap(err, "It looks like the data is not exist")
		return
	}
	if filter != nil && len(filter) > 0 {
		qq, args = utils.GenerateFilters(filter)
		where = fmt.Sprintf(" WHERE %s", qq)
		qq = fmt.Sprintf(" AND %s", qq)
	}
	q = fmt.Sprintf(`
		SELECT distribution_id, count(id) as subtotal 
		FROM %s 
		%s
		GROUP BY distribution_id
		ORDER BY distribution_id DESC
		LIMIT %d OFFSET %d`,
		tableFeedback, where, limit, (page-1)*limit,
	)
	items = make([]*Feedback, 0)
	rows, err = m.dbm.Query(ctx, m.dbm.Rebind(ctx, q), args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return items, 0, nil
		}
		return
	}
	distIds := make([]int64, 0)
	for rows.Next() {
		var distId, subtotal sql.NullInt64
		if err = rows.Scan(&distId, &subtotal); err != nil {
			return
		}
		if !utils.InArray(utils.ArrayInt64(distIds).ToArrayInterface(), distId.Int64) {
			distIds = append(distIds, distId.Int64)
		}
	}
	if len(distIds) < 1 {
		err = errors.New("nothing found")
		total = 0
		return
	}
	q = fmt.Sprintf(
		`SELECT 
						id, distribution_id, distribution_topic, distribution_object_id, range_start, range_end,
						respondent_id, respondent_username, respondent_name, respondent_email, 
						respondent_group_id, respondent_group_name, respondent_org_id, respondent_org_name,
						respondent_role, respondent_assignment,
						recipient_id, recipient_username, recipient_name, recipient_email, 
						recipient_group_id, recipient_group_name, recipient_org_id, recipient_org_name,
						recipient_role, recipient_assignment,
						link_id, hash, status, created_at, updated_at
					FROM %s
					WHERE distribution_id IN (%s) %s
					ORDER BY distribution_id desc, recipient_name asc, created_at, updated_at desc
			`,
		tableFeedback,
		strings.TrimRight(strings.Repeat("?,", len(distIds)), ","),
		qq,
	)
	args = append(utils.ArrayInt64(distIds).ToArrayInterface(), args...)
	rows, err = m.dbm.Query(ctx, m.dbm.Rebind(ctx, q), args...)
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
			&item.RespondentRole,
			&item.RespondentAssignment,
			&item.RecipientId,
			&item.RecipientUsername,
			&item.RecipientName,
			&item.RecipientEmail,
			&item.RecipientGroupId,
			&item.RecipientGroupName,
			&item.RecipientOrgId,
			&item.RecipientOrgName,
			&item.RecipientRole,
			&item.RecipientAssignment,
			&item.LinkId,
			&item.Hash,
			&item.Status,
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
				?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?, 
				?, ?, 
				NOW()
			)`,
		)
		values = append(values,
			&item.DistributionId, &item.DistributionTopic, &item.DistributionObjectId, &item.RangeStart, &item.RangeEnd,
			&item.RespondentId, &item.RespondentUsername, &item.RespondentName, &item.RespondentEmail,
			&item.RespondentGroupId, &item.RespondentGroupName, &item.RespondentOrgId, &item.RespondentOrgName,
			&item.RespondentRole, &item.RespondentAssignment,
			&item.RecipientId, &item.RecipientUsername, &item.RecipientName, &item.RecipientEmail,
			&item.RecipientGroupId, &item.RecipientGroupName, &item.RecipientOrgId, &item.RecipientOrgName,
			&item.RecipientRole, &item.RecipientAssignment,
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

func (m *manager) UpsertMultiple(ctx context.Context, items []*Feedback) (successItems [][]int64, err error) {
	q := fmt.Sprintf(`
		INSERT INTO %s (
			distribution_id, distribution_topic, distribution_object_id, range_start, range_end,
			respondent_id, respondent_username, respondent_name, respondent_email, 
			respondent_group_id, respondent_group_name, respondent_org_id, respondent_org_name,
			respondent_role, respondent_assignment,
			recipient_id, recipient_username, recipient_name, recipient_email, 
			recipient_group_id, recipient_group_name, recipient_org_id, recipient_org_name,
			recipient_role, recipient_assignment,
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
				?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?,
				?, ?, ?, ?,
				NOW()
			)`,
		)
		values = append(values,
			&item.DistributionId, &item.DistributionTopic, &item.DistributionObjectId, &item.RangeStart, &item.RangeEnd,
			&item.RespondentId, &item.RespondentUsername, &item.RespondentName, &item.RespondentEmail,
			&item.RespondentGroupId, &item.RespondentGroupName, &item.RespondentOrgId, &item.RespondentOrgName,
			&item.RespondentRole, &item.RespondentAssignment,
			&item.RecipientId, &item.RecipientUsername, &item.RecipientName, &item.RecipientEmail,
			&item.RecipientGroupId, &item.RecipientGroupName, &item.RecipientOrgId, &item.RecipientOrgName,
			&item.RecipientRole, &item.RecipientAssignment,
			&item.LinkId, &item.Hash, &item.Status, &item.Content,
		)
	}
	q = m.dbm.Rebind(ctx, fmt.Sprintf(
		`%s %s 
			ON CONFLICT (
				distribution_id, distribution_object_id, respondent_id, recipient_id
			)
			DO 
			UPDATE SET 
				distribution_topic = EXCLUDED.distribution_topic,
				respondent_role = EXCLUDED.respondent_role,
				respondent_assignment = EXCLUDED.respondent_assignment,
				recipient_role = EXCLUDED.recipient_role,
				recipient_assignment = EXCLUDED.recipient_assignment,
				range_start = EXCLUDED.range_start,
				range_end = EXCLUDED.range_end,
				content = EXCLUDED.content,
				updated_at = NOW()
			WHERE %s."status" = 'not-started'::feedback_status
			RETURNING id, distribution_id, distribution_object_id, respondent_id, recipient_id;
		`, q, strings.Join(placeholders, ","), tableFeedback))
	var rows database.Rows
	rows, err = m.dbm.Query(ctx, q, values...)
	if err != nil {
		return
	}
	successItems = make([][]int64, 0)
	for rows.Next() {
		var id, distributionId, distObjectId, respondentId, recipientId sql.NullInt64
		err = rows.Scan(&id, &distributionId, &distObjectId, &respondentId, &recipientId)
		if err != nil {
			return
		}
		successItems = append(successItems, []int64{
			id.Int64, distributionId.Int64, distObjectId.Int64, respondentId.Int64, recipientId.Int64,
		})
	}
	if len(successItems) < 1 {
		return successItems, fmt.Errorf("no rows created")
	}
	return
}

func (m *manager) Update(ctx context.Context, item Feedback) error {
	if item.Id < 1 {
		return fmt.Errorf("please provide the correct identifier")
	}
	args := []interface{}{
		&item.DistributionId, &item.DistributionTopic, &item.DistributionObjectId, &item.RangeStart, &item.RangeEnd,
		&item.RespondentId, &item.RespondentUsername, &item.RespondentName, &item.RespondentEmail,
		&item.RespondentGroupId, &item.RespondentGroupName, &item.RespondentOrgId, &item.RespondentOrgName,
		&item.RespondentRole, &item.RespondentAssignment,
		&item.RecipientId, &item.RecipientUsername, &item.RecipientName, &item.RecipientEmail,
		&item.RecipientGroupId, &item.RecipientGroupName, &item.RecipientOrgId, &item.RecipientOrgName,
		&item.RecipientRole, &item.RecipientAssignment,
		&item.LinkId, &item.Hash, &item.Status, &item.Content,
	}
	args = append(args, item.Id)
	q := fmt.Sprintf(`
		UPDATE %s 
		SET 
			distribution_id = ?, distribution_topic = ?, distribution_object_id = ?, 
			range_start = ?, range_end = ?,
			respondent_id = ?,  respondent_username = ?, respondent_name = ?, respondent_email = ?, 
			respondent_group_id = ?,  respondent_group_name = ?, respondent_org_id = ?, respondent_org_name = ?,
			respondent_role = ?, respondent_assignment = ?,
			recipient_id = ?, recipient_username = ?, recipient_name = ?, recipient_email = ?, 
			recipient_group_id = ?, recipient_group_name = ?, recipient_org_id = ?, recipient_org_name = ?,
			recipient_role = ?, recipient_assignment = ?,
			link_id = ?, hash = ?, status = ?, content = ?, updated_at = NOW()
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
