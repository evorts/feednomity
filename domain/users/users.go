package users

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

type IUsers interface {
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindGroupByIds(ctx context.Context, ids ...int64) ([]*Group, error)
	FindAll(ctx context.Context, page, limit int) (u []*User, total int, err error)
}

const (
	tableUsers      = "users"
	tableUsersGroup = "users_group"
)

func NewUserDomain(dbm database.IManager) IUsers {
	return &manager{dbm: dbm}
}

func (m *manager) FindByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := m.dbm.QueryRowAndBind(ctx, fmt.Sprintf(`
		SELECT 
			id, username, display_name, email, phone, password, role, group_id,
			created_at, updated_at 
		FROM %s WHERE username = $1
	`, tableUsers), []interface{}{username}, &user)

	if err != nil {
		return nil, errors.WithMessage(err, "fail to query user")
	}
	return &user, nil
}

func (m *manager) FindAll(ctx context.Context, page, limit int) (u []*User, total int, err error) {
	var (
		rows database.Rows
	)
	q := fmt.Sprintf(`SELECT count(id) FROM %s`, tableUsers)
	u = make([]*User, 0)
	err = m.dbm.QueryRowAndBind(ctx, q, nil, &total)
	if err != nil || total < 1 {
		err = errors.Wrap(err, "It looks like the data is not exist")
		return
	}
	rows, err = m.dbm.Query(
		ctx, fmt.Sprintf(
			`SELECT 
						id, username, display_name, email, phone, password, role, group_id,
						created_at, updated_at
					FROM %s LIMIT %d OFFSET %d`,
			tableUsers, limit, (page-1)*limit),
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, total, nil
		}
		return
	}
	for rows.Next() {
		var ui User
		err = rows.Scan(
			&ui.Id,
			&ui.Username,
			&ui.DisplayName,
			&ui.Email,
			&ui.Phone,
			&ui.Password,
			&ui.Role,
			&ui.GroupId,
			&ui.CreatedDate,
			&ui.UpdatedDate,
		)
		if err != nil {
			return
		}
		u = append(u, &ui)
	}
	return
}

func (m *manager) FindGroupByIds(ctx context.Context, ids ...int64) ([]*Group, error) {
	var (
		err  error
		rows database.Rows
		g    = make([]*Group, 0)
	)
	q := m.dbm.Rebind(ctx, fmt.Sprintf(
		`SELECT 
						id, name, disabled,
						created_at, updated_at, disabled_at
					FROM %s WHERE id IN (%s)`,
		tableUsersGroup, strings.TrimRight(strings.Repeat("?,", len(ids)), ",")),
	)
	rows, err = m.dbm.Query(
		ctx, q, utils.ArrayInt64(ids).ToArrayInterface()...,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return g, nil
		}
		return nil, err
	}
	for rows.Next() {
		var ug Group
		err = rows.Scan(
			&ug.Id,
			&ug.Name,
			&ug.Disabled,
			&ug.CreatedAt,
			&ug.UpdatedAt,
			&ug.DisabledAt,
		)
		if err != nil {
			return nil, err
		}
		g = append(g, &ug)
	}
	return g, nil
}
