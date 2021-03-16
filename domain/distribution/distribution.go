package distribution

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/utils"
	"strings"
)

type IManager interface {
	FindByIds(ctx context.Context, ids ...int64) ([]*Distribution, error)
	FindObjectByIds(ctx context.Context, ids ...int64) ([]*Object, error)
}

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
