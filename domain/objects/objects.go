package objects

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

type IObjects interface {
	FindAll(ctx context.Context, page, limit int) ([]*Object, int, error)
	FindByIds(ctx context.Context, id ...int64) ([]*Object, error)
	FindByName(ctx context.Context, name string, userGroupId int) (*Object, error)

	Insert(ctx context.Context, object Object) error
	DeleteByIds(ctx context.Context, id []int) error
}

const (
	tableObjects = "objects"
)

func NewObjectDomain(dbm database.IManager) IObjects {
	return &manager{dbm: dbm}
}

func (m *manager) FindAll(ctx context.Context, page, limit int) (o []*Object, total int, err error) {
	var (
		rows database.Rows
	)
	q := fmt.Sprintf(`SELECT count(id) FROM %s`, tableObjects)
	o = make([]*Object, 0)
	err = m.dbm.QueryRowAndBind(ctx, q, nil, &total)
	if err != nil || total < 1 {
		err = errors.Wrap(err, "It looks like the data is not exist")
		return
	}
	rows, err = m.dbm.Query(
		ctx, fmt.Sprintf(
			`SELECT 
						id, name, attributes, email, phone, role, assignment, user_group_id,
						disabled, archived, created_at, updated_at, disabled_at, archived_at
					FROM %s LIMIT %d OFFSET %d`,
			tableObjects, limit, (page-1)*limit), nil,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return o, total, nil
		}
		return
	}
	for rows.Next() {
		var object Object
		err = rows.Scan(
			&object.Id,
			&object.Name,
			&object.Attributes,
			&object.Email,
			&object.Phone,
			&object.Role,
			&object.Assignment,
			&object.UserGroupId,
			&object.Disabled,
			&object.Archived,
			&object.CreatedAt,
			&object.UpdatedAt,
			&object.DisabledAt,
			&object.ArchivedAt,
		)
		if err != nil {
			return
		}
		o = append(o, &object)
	}
	return
}

func (m *manager) FindByIds(ctx context.Context, ids ...int64) (o []*Object, err error) {
	o = make([]*Object, 0)
	q := m.dbm.Rebind(ctx,fmt.Sprintf(
		`SELECT
					id, name, attributes, email, phone, role, assignment, user_group_id,
						disabled, archived, created_at, updated_at, disabled_at, archived_at
				FROM %s WHERE id IN (%s)`, tableObjects,
		strings.TrimRight(strings.Repeat("?,", len(ids)), ","),
	))
	rows, err := m.dbm.Query(
		ctx, q,
		utils.ArrayInt64(ids).ToArrayInterface()...,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return o, nil
		}
		return
	}
	for rows.Next() {
		var (
			object Object
			phone sql.NullString
		)
		err = rows.Scan(
			&object.Id,
			&object.Name,
			&object.Attributes,
			&object.Email,
			&phone,
			&object.Role,
			&object.Assignment,
			&object.UserGroupId,
			&object.Disabled,
			&object.Archived,
			&object.CreatedAt,
			&object.UpdatedAt,
			&object.DisabledAt,
			&object.ArchivedAt,
		)
		object.Phone = phone.String
		if err != nil {
			return
		}
		o = append(o, &object)
	}
	return o, err
}

func (m *manager) FindByName(ctx context.Context, name string, userGroupId int) (*Object, error) {
	var object *Object
	err := m.dbm.QueryRowAndBind(
		ctx, fmt.Sprintf(`SELECT * FROM %s WHERE name=$1 AND user_group_id=$2`, tableObjects),
		[]interface{}{name, userGroupId}, object,
	)
	return object, err
}

func (m *manager) Insert(ctx context.Context, object Object) error {
	if object.Attributes == nil {
		object.Attributes = make(map[string]interface{}, 0)
	}
	q := m.dbm.Rebind(ctx, fmt.Sprintf(`
			INSERT INTO %s (name, attributes, email, phone, role, assignment, user_group_id, created_at)
			VALUES(?, ?, ?, ?, ?, ?, ?, NOW())
		`, tableObjects))
	rs, err := m.dbm.Exec(
		ctx, q,
		object.Name, object.Attributes, object.Email, object.Phone, object.Role, object.Assignment, object.UserGroupId,
	)
	if err != nil {
		return err
	}
	if rs.RowsAffected() < 1 {
		return errors.New("unable to insert new record correctly")
	}
	return nil
}

func (m *manager) DeleteByIds(ctx context.Context, ids []int) error {
	rs, err := m.dbm.Exec(
		ctx,
		fmt.Sprintf(`
			DELETE FROM %s WHERE id IN (%s)
		`, tableObjects, strings.TrimRight(strings.Repeat("?,", len(ids)), ",")),
		utils.ArrayInteger(ids).ToArrayInterface()...,
	)
	if err != nil {
		return err
	}
	if rs.RowsAffected() < 1 {
		return errors.New("not a single record removed")
	}
	return nil
}
