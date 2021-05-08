package acl

import (
	"context"
	"fmt"
	"github.com/evorts/feednomity/domain/users"
)

type AccessScope string

const (
	AccessScopeSelf   AccessScope = "self"
	AccessScopeGroup  AccessScope = "group"
	AccessScopeOrg    AccessScope = "org"
	AccessScopeGlobal AccessScope = "global"
)

type IManager interface {
	Populate() error
	IsAllowed(userId int64, method, path string) (allowed bool, scope AccessScope)
}

type access struct {
	Path             string      `json:"path"`
	Regex            bool        `json:"regex"`
	MethodAllowed    []string    `json:"method_allowed"`
	MethodDisallowed []string    `json:"method_disallowed"`
	AccessScope      AccessScope `json:"access_scope"`
	Disabled         bool        `json:"disabled"`
	AccessLevel      string      `json:"access_level"`
}

type accessControl struct {
	UserId     int      `json:"user_id"`
	Email      string   `json:"email"`
	GroupId    int      `json:"group_id"`
	Role       string   `json:"role"`
	RoleAccess []access `json:"role_access"`
	UserAccess []access `json:"user_access"`
}

type manager struct {
	u  users.IUsers
	ua users.IUserAccess
	ac map[string]accessControl // user_ud => accessControl
}

func NewACLManager(u users.IUsers, ua users.IUserAccess) IManager {
	return &manager{u: u, ua: ua}
}

func (m *manager) IsAllowed(userId int64, method, path string) (allowed bool, scope AccessScope) {
	fmt.Println(userId, method, path)
	allowed = true
	scope = AccessScopeGlobal
	return
}

func (m *manager) Populate() error {
	//read users data
	u, uErr := m.recursiveFindUsers(context.TODO(), 1, 10)
	if uErr != nil {
		return uErr
	}
	//read user role
	ur, urErr := m.recursiveFindRoleAccess(context.TODO(), 1, 10)
	if urErr != nil {
		return urErr
	}
	//read user access
	ua, uaErr := m.recursiveFindUserAccess(context.TODO(), 1, 10)
	if uaErr != nil {
		return uaErr
	}
	//populate to access control
	fmt.Println(u, ur, ua)
	return nil
}

func (m *manager) recursiveFindUsers(ctx context.Context, page, limit int) ([]*users.User, error) {
	u, ut, uErr := m.u.FindAll(ctx, page, limit)
	if uErr != nil {
		return nil, uErr
	}
	if (page-1)*limit > ut {
		return u, nil
	}
	uu, err := m.recursiveFindUsers(ctx, page+1, limit)
	if err != nil {
		return u, err
	}
	return append(u, uu...), nil
}

func (m *manager) recursiveFindRoleAccess(ctx context.Context, page, limit int) ([]*users.UserRoleAccess, error) {
	ur, uRoleTotal, uRoleErr := m.ua.FindAllRoleAccess(ctx, page, limit)
	if uRoleErr != nil {
		return nil, uRoleErr
	}
	if (page-1)*limit > uRoleTotal {
		return ur, nil
	}
	uu, err := m.recursiveFindRoleAccess(ctx, page+1, limit)
	if err != nil {
		return ur, err
	}
	return append(ur, uu...), nil
}

func (m *manager) recursiveFindUserAccess(ctx context.Context, page, limit int) ([]*users.UserAccess, error) {
	ua, uat, uaErr := m.ua.FindAllUserAccess(ctx, page, limit)
	if uaErr != nil {
		return nil, uaErr
	}
	if (page-1)*limit > uat {
		return ua, nil
	}
	uu, err := m.recursiveFindUserAccess(ctx, page+1, limit)
	if err != nil {
		return ua, err
	}
	return append(ua, uu...), nil
}
