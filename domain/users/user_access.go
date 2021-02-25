package users

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/pkg/errors"
)

const (
	tableUserAccess = "user_access"
	tableRoleAccess = "role_access"
)

type accessManager struct {
	dbm database.IManager
}

type IUserAccess interface {
	FindAllRoleAccess(ctx context.Context, page, limit int) (access []*UserRoleAccess, total int, err error)
	FindAllUserAccess(ctx context.Context, page, limit int) (access []*UserAccess, total int, err error)
}

func NewUserAccessDomain(dbm database.IManager) IUserAccess {
	return &accessManager{dbm: dbm}
}

func (m *accessManager) FindAllRoleAccess(ctx context.Context, page, limit int) (access []*UserRoleAccess, total int, err error) {
	var (
		rows database.Rows
	)
	q := fmt.Sprintf(`SELECT count(id) FROM %s`, tableRoleAccess)
	access = make([]*UserRoleAccess, 0)
	err = m.dbm.QueryRowAndBind(ctx, q, nil, &total)
	if err != nil || total < 1 {
		err = errors.Wrap(err, "It looks like the data is not exist")
		return
	}
	rows, err = m.dbm.Query(
		ctx, fmt.Sprintf(
			`SELECT 
						id, role, path, access_allowed, disabled
					FROM %s LIMIT %d OFFSET %d`,
			tableRoleAccess, limit, (page-1)*limit), nil,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return access, total, nil
		}
		return
	}
	for rows.Next() {
		var ra UserRoleAccess
		err = rows.Scan(
			&ra.Id,
			&ra.Role,
			&ra.Path,
			&ra.AccessAllowed,
			&ra.Disabled,
		)
		if err != nil {
			return
		}
		access = append(access, &ra)
	}
	return
}

func (m *accessManager) FindAllUserAccess(ctx context.Context, page, limit int) (access []*UserAccess, total int, err error) {
	var (
		rows database.Rows
	)
	q := fmt.Sprintf(`SELECT count(id) FROM %s`, tableRoleAccess)
	access = make([]*UserAccess, 0)
	err = m.dbm.QueryRowAndBind(ctx, q, nil, &total)
	if err != nil || total < 1 {
		err = errors.Wrap(err, "It looks like the data is not exist")
		return
	}
	rows, err = m.dbm.Query(
		ctx, fmt.Sprintf(
			`SELECT 
						id, user_id, path, access_allowed, access_disallowed, disabled
					FROM %s LIMIT %d OFFSET %d`,
			tableUserAccess, limit, (page-1)*limit), nil,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return access, total, nil
		}
		return
	}
	for rows.Next() {
		var ra UserAccess
		err = rows.Scan(
			&ra.Id,
			&ra.UserId,
			&ra.Path,
			&ra.AccessAllowed,
			&ra.AccessDisallowed,
			&ra.Disabled,
		)
		if err != nil {
			return
		}
		access = append(access, &ra)
	}
	return
}
