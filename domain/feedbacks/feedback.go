package feedbacks

import (
	"context"
	"fmt"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/pkg/errors"
)

const (
	tableFeedback = "feedbacks"
	tableFeedbackDetail = "feedback_detail"
	tableFeedbackLog = "feedback_log"
)
type manager struct {
	dbm database.IManager
}

type IFeedback interface {
	Save(ctx context.Context, f Feedback) error
	SaveDetail(ctx context.Context, fd Detail) error
	SaveTx(ctx context.Context, f Feedback, fd Detail) error
	RemoveByIds(ctx context.Context, ids ...int) error

	FindById(ctx context.Context, id int) (*Feedback, error)
	FindByDistId(ctx context.Context, distId, distObjId int64) (*Feedback, error)
	FindDetailByHash(ctx context.Context, linkHash string) (*Detail, error)
}

func NewFeedbackDomain(dbm database.IManager) IFeedback {
	return &manager{dbm: dbm}
}

func (m *manager) Save(ctx context.Context, f Feedback) error {
	// insert when id not exist
	if f.Id < 1 {
		q := m.dbm.Rebind(ctx, fmt.Sprintf(`
			INSERT INTO %s (
				distribution_id, distribution_object_id, distribution_topic, user_group_id, user_group_name,
				user_id, user_name, user_display_name, disabled, created_at
			)
			VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
		`, tableFeedback))
		rs, err := m.dbm.Exec(
			ctx, q,
			f.DistributionId, f.DistributionObjectId, f.DistributionTopic, f.UserGroupId, f.UserGroupName,
			f.UserId, f.UserName, f.UserDisplayName, f.Disabled,
		)
		if err != nil {
			return err
		}
		if rs.RowsAffected() < 1 {
			return errors.New("unable to insert new record correctly")
		}
		return nil
	}
	q := m.dbm.Rebind(ctx, fmt.Sprintf(`
			UPDATE %s 
			SET
				distribution_id = ?, distribution_object_id = ?, distribution_topic = ?, 
				user_group_id = ?, user_group_name = ?,
				user_id = ?, user_name = ?, user_display_name = ?, 
				disabled = ?, updated_at = NOW(), disabled_at = ?
			WHERE id = ?
		`, tableFeedback))
	rs, err := m.dbm.Exec(
		ctx, q,
		f.DistributionId, f.DistributionObjectId, f.DistributionTopic,
		f.UserGroupId, f.UserGroupName,
		f.UserId, f.UserName, f.UserDisplayName,
		f.Disabled, f.DisabledAt,
	)
	if err != nil {
		return err
	}
	if rs.RowsAffected() < 1 {
		return errors.New("unable to update the record correctly")
	}
	return nil
}

func (m *manager) SaveDetail(ctx context.Context, fd Detail) error {
	// insert when id not exist
	if fd.Id < 1 {
		return nil
	}
	return nil
}

func (m *manager) SaveTx(ctx context.Context, f Feedback, fd Detail) error {
	var (
		ct database.CommandTag
		err error
		tx database.Tx
		actionF, actionFd = "insert", "insert"
		fid int64
	)
	tx, err = m.dbm.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	// feedback
	qf := m.dbm.Rebind(ctx, fmt.Sprintf(`
			INSERT INTO %s (
				distribution_id, distribution_object_id, distribution_topic, user_group_id, user_group_name,
				user_id, user_name, user_display_name, disabled, created_at, disabled_at
			)
			VALUES(
				?, ?, ?, ?, ?, 
				?, ?, ?, ?, NOW(), ?
			)
			RETURNING id
		`, tableFeedback))
	if f.Id > 0 {
		actionF = "update"
		qf = m.dbm.Rebind(ctx, fmt.Sprintf(`
			UPDATE %s 
			SET 
				distribution_id = ?, distribution_object_id = ?, distribution_topic = ?, 
				user_group_id = ?, user_group_name = ?, user_id = ?, user_name = ?, user_display_name = ?, disabled = ?, 
				updated_at = NOW(), disabled_at = ?
			WHERE id = %d
			RETURNING id
		`, tableFeedback, f.Id))
	}
	err = m.dbm.QueryRow(
		ctx, qf,
		f.DistributionId, f.DistributionObjectId, f.DistributionTopic,
		f.UserGroupId, f.UserGroupName, f.UserId, f.UserName, f.UserDisplayName, f.Disabled, f.DisabledAt,
	).Scan(&fid)
	if err != nil {
		return errors.New("failed to save feedback")
	}
	// feedback detail
	qfd :=  m.dbm.Rebind(ctx, fmt.Sprintf(`
			INSERT INTO %s (
				feedback_id, link_id, hash, respondent_id, respondent_name, respondent_email,
				recipient_id, recipient_name, recipient_email, content, status, created_at
			)
			VALUES(
				?, ?, ?, ?, ?, ?, 
				?, ?, ?, ?, ?, NOW()
			)
		`, tableFeedbackDetail))
	if fd.Id > 0 {
		actionFd = "update"
		qfd =  m.dbm.Rebind(ctx, fmt.Sprintf(`
			UPDATE %s 
			SET 
				feedback_id = ?, link_id = ?, hash = ?, respondent_id = ?, respondent_name = ?, respondent_email = ?,
				recipient_id = ?, recipient_name = ?, recipient_email = ?, content = ?, status = ?,  updated_at = NOW()
			WHERE id = %d
		`, tableFeedbackDetail, fd.Id))
	}
	ct, err = m.dbm.Exec(
		ctx, qfd,
		fid, fd.LinkId, fd.Hash, fd.RespondentId, fd.RespondentName, fd.RespondentEmail,
		fd.RecipientId, fd.RecipientName, fd.RecipientEmail, fd.Content, fd.Status,
	)
	if err != nil || ct.RowsAffected() < 1 {
		return errors.New("failed to save feedback detail")
	}
	fl := m.dbm.Rebind(ctx, fmt.Sprintf(`
			INSERT INTO %s (
				feedback_id, action, values, values_prev, notes, at
			)
			VALUES(
				?, ?, ?, ?, ?, NOW()
			)
		`, tableFeedbackLog))
	_, _ = m.dbm.Exec(
		ctx, fl,
		fid, fmt.Sprintf("save:%s", actionF), f, nil, "",
	)
	_, _ = m.dbm.Exec(
		ctx, fl,
		fid, fmt.Sprintf("save:%s", actionFd), fd, nil, "",
	)
	return tx.Commit(ctx)
}

func (m *manager) RemoveByIds(ctx context.Context, ids ...int) error {
	panic("implement me")
}

func (m *manager) FindById(ctx context.Context, id int) (*Feedback, error) {
	panic("implement me")
}

func (m *manager) FindByDistId(ctx context.Context, distId, distObjId int64) (*Feedback, error) {
	q := fmt.Sprintf(`
		SELECT 
			id, distribution_id, distribution_object_id, distribution_topic, user_group_id, user_group_name,
			user_id, user_name, user_display_name, disabled, created_at, updated_at, disabled_at
		FROM %s
		WHERE distribution_id = $1 and distribution_object_id = $2`, tableFeedback)
	var (
		f Feedback
	)
	err := m.dbm.QueryRowAndBind(ctx, q, []interface{}{distId, distObjId},
		&f.Id, &f.DistributionId, &f.DistributionObjectId, &f.DistributionTopic, &f.UserGroupId, &f.UserGroupName,
		&f.UserId, &f.UserName, &f.UserDisplayName, &f.Disabled, &f.CreatedAt, &f.UpdateAt, &f.DisabledAt,
	)
	return &f, err
}

func (m *manager) FindDetailByHash(ctx context.Context, linkHash string) (*Detail, error) {
	q := fmt.Sprintf(`
		SELECT 
			id, feedback_id, link_id, hash, respondent_id, respondent_name, respondent_email,
			recipient_id, recipient_name, recipient_email, content, status, created_at, updated_at
		FROM %s
		WHERE hash = $1`, tableFeedbackDetail)
	var (
		d Detail
	)
	err := m.dbm.QueryRowAndBind(ctx, q, []interface{}{linkHash},
		&d.Id, &d.FeedbackId, &d.LinkId, &d.Hash, &d.RespondentId, &d.RespondentName, &d.RespondentEmail,
		&d.RecipientId, &d.RecipientName, &d.RecipientEmail, &d.Content, &d.Status, &d.CreatedAt, &d.UpdatedAt,
	)
	return &d, err
}

