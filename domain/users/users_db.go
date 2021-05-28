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

const (
	tableUsers      = "users"
	tableUsersGroup = "users_group"
	tableUsersOrg   = "users_organization"
)

func NewUserDomain(dbm database.IManager) IUsers {
	return &manager{dbm: dbm}
}

func (m *manager) FindByIds(ctx context.Context, ids ...int64) ([]*User, error) {
	var (
		err  error
		rows database.Rows
		u    = make([]*User, 0)
	)
	q := m.dbm.Rebind(ctx, fmt.Sprintf(
		`SELECT 
						id, username, display_name, attributes, 
						email, phone, password, pin, 
						access_role, job_role, assignment, group_id,
						disabled, created_at, updated_at, disabled_at
					FROM %s WHERE id IN (%s)`,
		tableUsers, strings.TrimRight(strings.Repeat("?,", len(ids)), ",")),
	)
	rows, err = m.dbm.Query(
		ctx, q, utils.ArrayInt64(ids).ToArrayInterface()...,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return u, nil
		}
		return nil, err
	}
	for rows.Next() {
		var ui User
		var displayName, phone, pwd, pin, jobRole, assignment sql.NullString
		err = rows.Scan(
			&ui.Id,
			&ui.Username,
			&displayName,
			&ui.Attributes,
			&ui.Email,
			&phone,
			&pwd,
			&pin,
			&ui.AccessRole,
			&jobRole,
			&assignment,
			&ui.GroupId,
			&ui.Disabled,
			&ui.CreatedAt,
			&ui.UpdatedAt,
			&ui.DisabledAt,
		)
		ui.DisplayName = displayName.String
		ui.Phone = phone.String
		ui.Password = pwd.String
		ui.PIN = pin.String
		ui.JobRole = jobRole.String
		ui.Assignment = assignment.String
		if err != nil {
			return nil, err
		}
		u = append(u, &ui)
	}
	return u, nil
}

func (m *manager) FindByUserEmail(ctx context.Context, value string) (*User, error) {
	var item User
	var displayName, phone, pwd, pin, jobRole, assignment sql.NullString
	err := m.dbm.QueryRow(ctx, fmt.Sprintf(`
		SELECT 
			id, username, display_name, attributes, 
			email, phone, password, pin, 
			access_role, job_role, assignment, group_id,
			disabled, created_at, updated_at, disabled_at
		FROM %s WHERE email = $1
	`, tableUsers), value).Scan(
		&item.Id,
		&item.Username,
		&displayName,
		&item.Attributes,
		&item.Email,
		&phone,
		&pwd,
		&pin,
		&item.AccessRole,
		&jobRole,
		&assignment,
		&item.GroupId,
		&item.Disabled,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.DisabledAt,
	)
	if err != nil {
		return nil, errors.WithMessage(err, "fail to query user")
	}

	item.DisplayName = displayName.String
	item.Phone = phone.String
	item.Password = pwd.String
	item.PIN = pin.String
	item.JobRole = jobRole.String
	item.Assignment = assignment.String

	return &item, nil
}

func (m *manager) FindByUsername(ctx context.Context, username string) (*User, error) {
	var ui User
	var displayName, phone, pwd, pin, jobRole, assignment sql.NullString
	err := m.dbm.QueryRow(ctx, fmt.Sprintf(`
		SELECT 
			id, username, display_name, attributes, 
			email, phone, password, pin, 
			access_role, job_role, assignment, group_id,
			disabled, created_at, updated_at, disabled_at
		FROM %s WHERE username = $1
	`, tableUsers), username).Scan(
		&ui.Id,
		&ui.Username,
		&displayName,
		&ui.Attributes,
		&ui.Email,
		&phone,
		&pwd,
		&pin,
		&ui.AccessRole,
		&jobRole,
		&assignment,
		&ui.GroupId,
		&ui.Disabled,
		&ui.CreatedAt,
		&ui.UpdatedAt,
		&ui.DisabledAt,
	)
	if err != nil {
		return nil, errors.WithMessage(err, "fail to query user")
	}

	ui.DisplayName = displayName.String
	ui.Phone = phone.String
	ui.Password = pwd.String
	ui.PIN = pin.String
	ui.JobRole = jobRole.String
	ui.Assignment = assignment.String

	return &ui, nil
}

func (m *manager) FindByNameAndGroupId(ctx context.Context, name string, groupId int) ([]*User, error) {
	var u []*User
	q := m.dbm.Rebind(ctx, fmt.Sprintf(
		`SELECT 
						id, username, display_name, attributes, 
						email, phone,  
						access_role, job_role, assignment, group_id,
						disabled, created_at, updated_at, disabled_at
					FROM %s WHERE group_id=? AND display_name LIKE ?`,
		tableUsers),
	)
	rows, err := m.dbm.Query(ctx, q, groupId, "%"+name+"%")
	if err != nil {
		if err == sql.ErrNoRows {
			return u, nil
		}
		return u, err
	}
	for rows.Next() {
		var ui User
		var displayName, phone, jobRole, assignment sql.NullString
		err = rows.Scan(
			&ui.Id,
			&ui.Username,
			&displayName,
			&ui.Attributes,
			&ui.Email,
			&phone,
			&ui.AccessRole,
			&jobRole,
			&assignment,
			&ui.GroupId,
			&ui.Disabled,
			&ui.CreatedAt,
			&ui.UpdatedAt,
			&ui.DisabledAt,
		)
		ui.DisplayName = displayName.String
		ui.Phone = phone.String
		ui.JobRole = jobRole.String
		ui.Assignment = assignment.String
		if err != nil {
			return u, err
		}
		u = append(u, &ui)
	}
	return u, nil
}

func (m *manager) FindByNameAndOrgId(ctx context.Context, name string, orgId int) ([]*User, error) {
	var u []*User
	q := m.dbm.Rebind(ctx, fmt.Sprintf(
		`SELECT 
						id, username, display_name, attributes, 
						email, phone,  
						access_role, job_role, assignment, group_id,
						disabled, created_at, updated_at, disabled_at
					FROM %s WHERE org_id=? AND display_name LIKE ?`,
		tableUsers),
	)
	rows, err := m.dbm.Query(ctx, q, orgId, "%"+name+"%")
	if err != nil {
		if err == sql.ErrNoRows {
			return u, nil
		}
		return u, err
	}
	for rows.Next() {
		var ui User
		var displayName, phone, jobRole, assignment sql.NullString
		err = rows.Scan(
			&ui.Id,
			&ui.Username,
			&displayName,
			&ui.Attributes,
			&ui.Email,
			&phone,
			&ui.AccessRole,
			&jobRole,
			&assignment,
			&ui.GroupId,
			&ui.Disabled,
			&ui.CreatedAt,
			&ui.UpdatedAt,
			&ui.DisabledAt,
		)
		ui.DisplayName = displayName.String
		ui.Phone = phone.String
		ui.JobRole = jobRole.String
		ui.Assignment = assignment.String
		if err != nil {
			return u, err
		}
		u = append(u, &ui)
	}
	return u, nil
}

func (m *manager) FindByName(ctx context.Context, name string) ([]*User, error) {
	var u []*User
	q := m.dbm.Rebind(ctx, fmt.Sprintf(
		`SELECT 
						id, username, display_name, attributes, 
						email, phone,  
						access_role, job_role, assignment, group_id,
						disabled, created_at, updated_at, disabled_at
					FROM %s WHERE display_name LIKE ?`,
		tableUsers),
	)
	rows, err := m.dbm.Query(ctx, q, "%"+name+"%")
	if err != nil {
		if err == sql.ErrNoRows {
			return u, nil
		}
		return u, err
	}
	for rows.Next() {
		var ui User
		var displayName, phone, jobRole, assignment sql.NullString
		err = rows.Scan(
			&ui.Id,
			&ui.Username,
			&displayName,
			&ui.Attributes,
			&ui.Email,
			&phone,
			&ui.AccessRole,
			&jobRole,
			&assignment,
			&ui.GroupId,
			&ui.Disabled,
			&ui.CreatedAt,
			&ui.UpdatedAt,
			&ui.DisabledAt,
		)
		ui.DisplayName = displayName.String
		ui.Phone = phone.String
		ui.JobRole = jobRole.String
		ui.Assignment = assignment.String
		if err != nil {
			return u, err
		}
		u = append(u, &ui)
	}
	return u, nil
}

func (m *manager) FindAll(ctx context.Context, page, limit int) (items []*User, total int, err error) {
	var (
		rows database.Rows
	)
	q := fmt.Sprintf(`SELECT count(id) FROM %s`, tableUsers)
	items = make([]*User, 0)
	err = m.dbm.QueryRowAndBind(ctx, q, nil, &total)
	if err != nil || total < 1 {
		err = errors.Wrap(err, "It looks like the data is not exist")
		return
	}
	rows, err = m.dbm.Query(
		ctx, fmt.Sprintf(
			`SELECT 
						id, username, display_name, attributes, 
						email, phone,  
						access_role, job_role, assignment, group_id,
						disabled, created_at, updated_at, disabled_at
					FROM %s LIMIT %d OFFSET %d`,
			tableUsers, limit, (page-1)*limit),
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return items, total, nil
		}
		return
	}
	for rows.Next() {
		var ui User
		var phone, assignment sql.NullString
		err = rows.Scan(
			&ui.Id,
			&ui.Username,
			&ui.DisplayName,
			&ui.Attributes,
			&ui.Email,
			&phone,
			&ui.AccessRole,
			&ui.JobRole,
			&assignment,
			&ui.GroupId,
			&ui.Disabled,
			&ui.CreatedAt,
			&ui.UpdatedAt,
			&ui.DisabledAt,
		)
		ui.Phone = phone.String
		ui.Assignment = assignment.String
		if err != nil {
			return
		}
		items = append(items, &ui)
	}
	return
}

func (m *manager) Insert(ctx context.Context, item User) error {
	if item.Attributes == nil {
		item.Attributes = make(map[string]interface{}, 0)
	}
	var disabledAt, pwd, pin interface{} = nil, nil, nil
	var pwdArg, pinArg = "?", "?"
	if len(item.Password) > 0 {
		pwdArg = "digest(?, 'sha1')"
		pwd = item.Password
	}
	if len(item.PIN) > 0 {
		pinArg = "digest(?, 'sha1')"
		pin = item.PIN
	}
	q := m.dbm.Rebind(ctx, fmt.Sprintf(`
			INSERT INTO %s (
				username, display_name, attributes, email, phone,
				password, pin, access_role, job_role, assignment, group_id,
				disabled, created_at, disabled_at
			)
			VALUES (
				?, ?, ?, ?, ?, 
				%s, %s, ?, ?, ?, ?, 
				?, NOW(), ?
			)
		`, tableUsers, pwdArg, pinArg))
	if item.Disabled {
		disabledAt = "NOW()"
	}
	p, err := m.dbm.Prepare(ctx, "users_insert", q)
	if err != nil {
		return err
	}
	rs, err := m.dbm.Exec(
		ctx, p.SQL,
		item.Username, item.DisplayName, item.Attributes, item.Email, item.Phone,
		pwd, pin, item.AccessRole, item.JobRole, item.Assignment, item.GroupId,
		item.Disabled, disabledAt,
	)
	if err != nil {
		return err
	}
	if rs.RowsAffected() < 1 {
		return errors.New("unable to insert new record correctly")
	}
	return nil
}

func (m *manager) InsertMultiple(ctx context.Context, items []*User) error {
	q := fmt.Sprintf(`
		INSERT INTO %s (
				username, display_name, attributes, email, phone,
				password, pin, access_role, job_role, assignment, group_id,
				disabled, created_at, disabled_at)
			VALUES`, tableUsers)
	placeholders := make([]string, 0)
	values := make([]interface{}, 0)
	for _, usv := range items {
		var pwdArg, pinArg = "?", "?"
		var disabledAt, pwd, pin interface{} = nil, nil, nil
		if len(usv.Password) > 0 {
			pwdArg = "digest(?, 'sha1')"
			pwd = usv.Password
		}
		if len(usv.PIN) > 0 {
			pinArg = "digest(?, 'sha1')"
			pin = usv.PIN
		}
		placeholders = append(
			placeholders,
			fmt.Sprintf(`(
				?, ?, ?, ?, ?, 
				%s, %s, ?, ?, ?, ?, 
				?, NOW(), ?
			)`, pwdArg, pinArg),
		)
		if usv.Disabled {
			disabledAt = "NOW()"
		}
		values = append(
			values, usv.Username, usv.DisplayName, usv.Attributes, usv.Email, usv.Phone,
			pwd, pin, usv.AccessRole, usv.JobRole, usv.Assignment, usv.GroupId,
			usv.Disabled, disabledAt,
		)
	}
	q = m.dbm.Rebind(ctx, fmt.Sprintf(`%s %s`, q, strings.Join(placeholders, ",")))
	cmd, err2 := m.dbm.Exec(ctx, q, values...)
	if err2 != nil {
		return errors.Wrap(err2, "failed saving users. some errors in constraint or data.")
	}
	if cmd.RowsAffected() > 0 {
		return nil
	}
	return fmt.Errorf("no rows created")
}

func (m *manager) Update(ctx context.Context, item User) error {
	if item.Id < 1 {
		return fmt.Errorf("please provide the correct user identifier")
	}
	var (
		phone interface{} = nil
	)
	if len(item.Phone) > 0 {
		phone = item.Phone
	}
	args := []interface{}{
		item.Username,
		item.DisplayName,
		item.Attributes,
		item.Email,
		phone,
		item.AccessRole,
		item.JobRole,
		item.Assignment,
		item.GroupId,
		item.Disabled,
	}
	var pwdArg, pinArg = "", ""
	if len(item.Password) > 0 {
		pwdArg = "password = digest(?, 'sha1'),"
		args = append(args, item.Password)
	}
	if len(item.PIN) > 0 {
		pinArg = "pin = digest(?, 'sha1'),"
		args = append(args, item.PIN)
	}
	var disabledAt interface{} = item.DisabledAt
	if item.Disabled {
		disabledAt = "NOW()"
	}
	args = append(args, disabledAt)
	args = append(args, item.Id)
	q := fmt.Sprintf(`
		UPDATE %s 
		SET 
			username = ?,
			display_name = ?,
			attributes = ?,
			email = ?,
			phone = ?,
			access_role = ?,
			job_role = ?,
			assignment = ?,
			group_id = ?,
			disabled = ?,
			updated_at = NOW(),
			%s %s
			disabled_at = ?
		WHERE id = ?`, tableUsers, pwdArg, pinArg)
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

func (m *manager) UpdatePasswordById(ctx context.Context, id int64, pass string) error {
	q := m.dbm.Rebind(ctx, fmt.Sprintf(`UPDATE %s SET password = digest(?, 'sha1') WHERE id = ?`, tableUsers))
	cmd, err := m.dbm.Exec(ctx, q, pass, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() > 0 {
		return nil
	}
	return fmt.Errorf("no rows updated")
}

func (m *manager) DeleteByIds(ctx context.Context, id []int64) error {
	q := m.dbm.Rebind(ctx, fmt.Sprintf(`
			DELETE FROM %s WHERE id IN (%s)
		`, tableUsers, strings.TrimRight(strings.Repeat("?,", len(id)), ",")))
	rs, err := m.dbm.Exec(
		ctx,
		q,
		utils.ArrayInt64(id).ToArrayInterface()...,
	)
	if err != nil {
		return err
	}
	if rs.RowsAffected() < 1 {
		return errors.New("not a single record removed")
	}
	return nil
}

func (m *manager) DisableByIds(ctx context.Context, id []int64) error {
	q := m.dbm.Rebind(ctx, fmt.Sprintf(`
			UPDATE %s 
			SET disabled = true, disabled_at = NOW()
			WHERE id IN (%s)
		`, tableUsers, strings.TrimRight(strings.Repeat("?,", len(id)), ",")))

	rs, err := m.dbm.Exec(ctx, q, utils.ArrayInt64(id).ToArrayInterface()...)
	if err != nil {
		return err
	}
	if rs.RowsAffected() < 1 {
		return errors.New("not a single record removed")
	}
	return nil
}

func (m *manager) FindGroupByIds(ctx context.Context, ids ...int64) ([]*Group, error) {
	var (
		err  error
		rows database.Rows
		g    = make([]*Group, 0)
	)
	q := m.dbm.Rebind(ctx, fmt.Sprintf(
		`SELECT 
						id, name, org_id, disabled,
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
			&ug.OrgId,
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

func (m *manager) FindGroupByOrgId(ctx context.Context, id int64) ([]*Group, error) {
	var (
		err  error
		rows database.Rows
		g    = make([]*Group, 0)
	)
	q := fmt.Sprintf(
		`SELECT 
					id, name, org_id, disabled,
					created_at, updated_at, disabled_at
				FROM %s WHERE org_id = $1`,
		tableUsersGroup,
	)
	rows, err = m.dbm.Query(
		ctx, q, id,
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
			&ug.OrgId,
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

func (m *manager) FindAllGroups(ctx context.Context, page, limit int) (items []*Group, total int, err error) {
	var (
		rows database.Rows
	)
	q := fmt.Sprintf(`SELECT count(id) FROM %s`, tableUsersGroup)
	items = make([]*Group, 0)
	err = m.dbm.QueryRowAndBind(ctx, q, nil, &total)
	if err != nil || total < 1 {
		err = errors.Wrap(err, "It looks like the data is not exist")
		return
	}
	rows, err = m.dbm.Query(
		ctx, fmt.Sprintf(
			`SELECT 
						id, name, org_id, disabled, created_at, updated_at, disabled_at
					FROM %s LIMIT %d OFFSET %d`,
			tableUsersGroup, limit, (page-1)*limit),
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return items, total, nil
		}
		return
	}
	for rows.Next() {
		var item Group
		err = rows.Scan(
			&item.Id,
			&item.Name,
			&item.OrgId,
			&item.Disabled,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.DisabledAt,
		)
		if err != nil {
			return
		}
		items = append(items, &item)
	}
	return
}

func (m *manager) InsertGroup(ctx context.Context, item Group) error {
	var disabledAt interface{} = nil
	q := m.dbm.Rebind(ctx, fmt.Sprintf(`
			INSERT INTO %s (
				name, org_id, disabled, created_at, disabled_at
			)
			VALUES (
				?, ?, ?, NOW(), ?
			)
		`, tableUsersGroup))
	if item.Disabled {
		disabledAt = "NOW()"
	}
	p, err := m.dbm.Prepare(ctx, "groups_insert", q)
	if err != nil {
		return err
	}
	rs, err := m.dbm.Exec(
		ctx, p.SQL,
		item.Name, item.OrgId, item.Disabled, disabledAt,
	)
	if err != nil {
		return err
	}
	if rs.RowsAffected() < 1 {
		return errors.New("unable to insert new record correctly")
	}
	return nil
}

func (m *manager) InsertGroups(ctx context.Context, items []*Group) error {
	q := fmt.Sprintf(`
		INSERT INTO %s (
				name, org_id, disabled, created_at, disabled_at)
			VALUES`, tableUsersGroup)
	placeholders := make([]string, 0)
	values := make([]interface{}, 0)
	for _, item := range items {
		var disabledAt interface{} = nil
		placeholders = append(
			placeholders,
			`(?, ?, ?, NOW(), ?)`,
		)
		if item.Disabled {
			disabledAt = "NOW()"
		}
		values = append(
			values, item.Name, item.OrgId, item.Disabled, disabledAt,
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

func (m *manager) UpdateGroup(ctx context.Context, item Group) error {
	if item.Id < 1 {
		return fmt.Errorf("please provide the correct identifier")
	}
	args := []interface{}{
		item.Name,
		item.OrgId,
		item.Disabled,
	}
	var disabledAt interface{} = item.DisabledAt
	if item.Disabled {
		disabledAt = "NOW()"
	}
	args = append(args, disabledAt)
	args = append(args, item.Id)
	q := fmt.Sprintf(`
		UPDATE %s 
		SET 
			name = ?,
			org_id = ?,
			disabled = ?,
			updated_at = NOW(),
			disabled_at = ?
		WHERE id = ?`, tableUsersGroup)
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

func (m *manager) DeleteGroupByIds(ctx context.Context, ids ...int64) error {
	q := m.dbm.Rebind(ctx, fmt.Sprintf(`
			DELETE FROM %s WHERE id IN (%s)
		`, tableUsersGroup, strings.TrimRight(strings.Repeat("?,", len(ids)), ",")))
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

func (m *manager) FindOrganizationByIds(ctx context.Context, ids ...int64) ([]*Organization, error) {
	var (
		err   error
		rows  database.Rows
		items = make([]*Organization, 0)
	)
	q := m.dbm.Rebind(ctx, fmt.Sprintf(
		`SELECT 
						id, name, address, phone, disabled,
						created_at, updated_at, disabled_at
					FROM %s WHERE id IN (%s)`,
		tableUsersOrg, strings.TrimRight(strings.Repeat("?,", len(ids)), ",")),
	)
	rows, err = m.dbm.Query(
		ctx, q, utils.ArrayInt64(ids).ToArrayInterface()...,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return items, nil
		}
		return nil, err
	}
	for rows.Next() {
		var (
			item           Organization
			phone, address sql.NullString
			disabled       sql.NullBool
		)
		err = rows.Scan(
			&item.Id,
			&item.Name,
			&address,
			&phone,
			&disabled,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.DisabledAt,
		)
		item.Address = address.String
		item.Phone = phone.String
		item.Disabled = disabled.Bool
		if err != nil {
			return items, err
		}
		items = append(items, &item)
	}
	return items, nil
}

func (m *manager) FindAllOrganizations(ctx context.Context, page, limit int) (items []*Organization, total int, err error) {
	var (
		rows database.Rows
	)
	q := fmt.Sprintf(`SELECT count(id) FROM %s`, tableUsersGroup)
	items = make([]*Organization, 0)
	err = m.dbm.QueryRowAndBind(ctx, q, nil, &total)
	if err != nil || total < 1 {
		err = errors.Wrap(err, "It looks like the data is not exist")
		return
	}
	rows, err = m.dbm.Query(
		ctx, fmt.Sprintf(
			`SELECT 
						id, name, address, phone, disabled, created_at, updated_at, disabled_at
					FROM %s LIMIT %d OFFSET %d`,
			tableUsersOrg, limit, (page-1)*limit),
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return items, total, nil
		}
		return
	}
	for rows.Next() {
		var (
			item     Organization
			phone    sql.NullString
			disabled sql.NullBool
		)
		err = rows.Scan(
			&item.Id,
			&item.Name,
			&item.Address,
			&phone,
			&disabled,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.DisabledAt,
		)
		item.Phone = phone.String
		item.Disabled = disabled.Bool
		if err != nil {
			return
		}
		items = append(items, &item)
	}
	return
}

func (m *manager) InsertOrganization(ctx context.Context, item Organization) error {
	var disabledAt interface{} = nil
	q := m.dbm.Rebind(ctx, fmt.Sprintf(`
			INSERT INTO %s (
				name, address, phone, disabled, created_at, disabled_at
			)
			VALUES (
				?, ?, ?, ?, NOW(), ?
			)
		`, tableUsersOrg))
	if item.Disabled {
		disabledAt = "NOW()"
	}
	p, err := m.dbm.Prepare(ctx, "organizations_insert", q)
	if err != nil {
		return err
	}
	rs, err := m.dbm.Exec(
		ctx, p.SQL,
		item.Name, item.Address, item.Phone, item.Disabled, disabledAt,
	)
	if err != nil {
		return err
	}
	if rs.RowsAffected() < 1 {
		return errors.New("unable to insert new record correctly")
	}
	return nil
}

func (m *manager) InsertOrganizations(ctx context.Context, items []*Organization) error {
	q := fmt.Sprintf(`
		INSERT INTO %s (
				name, address, phone, disabled, created_at, disabled_at)
			VALUES`, tableUsersOrg)
	placeholders := make([]string, 0)
	values := make([]interface{}, 0)
	for _, item := range items {
		var disabledAt interface{} = nil
		placeholders = append(
			placeholders,
			`(?, ?, ?, ?, NOW(), ?)`,
		)
		if item.Disabled {
			disabledAt = "NOW()"
		}
		values = append(
			values, item.Name, item.Address, item.Phone, item.Disabled, disabledAt,
		)
	}
	q = m.dbm.Rebind(ctx, fmt.Sprintf(`%s %s`, q, strings.Join(placeholders, ",")))
	cmd, err2 := m.dbm.Exec(ctx, q, values...)
	if err2 != nil {
		return errors.Wrap(err2, "failed saving organizations. some errors in constraint or data.")
	}
	if cmd.RowsAffected() > 0 {
		return nil
	}
	return fmt.Errorf("no rows created")
}

func (m *manager) UpdateOrganization(ctx context.Context, item Organization) error {
	if item.Id < 1 {
		return fmt.Errorf("please provide the correct identifier")
	}
	args := []interface{}{
		item.Name,
		item.Address,
		item.Phone,
		item.Disabled,
	}
	var disabledAt interface{} = item.DisabledAt
	if item.Disabled {
		disabledAt = "NOW()"
	}
	args = append(args, disabledAt)
	args = append(args, item.Id)
	q := fmt.Sprintf(`
		UPDATE %s 
		SET 
			name = ?,
			address = ?,
			phone = ?,
			disabled = ?,
			updated_at = NOW(),
			disabled_at = ?
		WHERE id = ?`, tableUsersOrg)
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

func (m *manager) DeleteOrganizationByIds(ctx context.Context, ids ...int64) error {
	q := m.dbm.Rebind(ctx, fmt.Sprintf(`
			DELETE FROM %s WHERE id IN (%s)
		`, tableUsersOrg, strings.TrimRight(strings.Repeat("?,", len(ids)), ",")))
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
