package users

import (
	"context"
	"fmt"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/pkg/errors"
)

type manager struct {
	dbm database.IManager
}

const (
	tableUsers      = "users"
	tableUserAccess = "user_access"
	tableRoleAccess = "role_access"
)

func (m *manager) FindByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := m.dbm.QueryRowAndBind(ctx, fmt.Sprintf(`
		SELECT * FROM %s WHERE username = $1
	`, tableUsers), []interface{}{username}, &user)

	if err != nil {
		return nil, errors.WithMessage(err, "fail to query user")
	}
	return &user, nil
}

type IUsers interface {
	FindByUsername(ctx context.Context, username string) (*User, error)
}

func NewUserDomain(dbm database.IManager) IUsers {
	return &manager{dbm: dbm}
}
