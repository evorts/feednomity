package acl

import (
	"context"
	"fmt"
	"github.com/evorts/feednomity/domain/users"
)

type IManager interface {
	Populate() error
	IsAllowed(userId int, path, method string) bool
}

type access struct {
	Path          string   `json:"path"`
	MethodAllowed []string `json:"method_allowed"`
	Disabled      bool     `json:"disabled"`
	AccessLevel   string   `json:"access_level"`
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

func (m *manager) IsAllowed(userId int, path, method string) bool {
	return true
}

func (m *manager) Populate() error {
	//read users data
	u, uerr := m.recursiveFindUsers(context.TODO(), 1, 10)
	if uerr != nil {
		return uerr
	}
	//read user role
	ur, urerr := m.recursiveFindRoleAccess(context.TODO(), 1, 10)
	if urerr != nil {
		return urerr
	}
	//read user access
	ua, uaerr := m.recursiveFindUserAccess(context.TODO(), 1, 10)
	if uaerr != nil {
		return uaerr
	}
	//populate to access control
	fmt.Println(u, ur, ua)
	return nil
}

func (m *manager) recursiveFindUsers(ctx context.Context, page, limit int) ([]*users.User, error) {
	u, ut, uerr := m.u.FindAll(ctx, page, limit)
	if uerr != nil {
		return nil, uerr
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
	ur, urt, urerr := m.ua.FindAllRoleAccess(ctx, 1, 10)
	if urerr != nil {
		return nil, urerr
	}
	if (page-1)*limit > urt {
		return ur, nil
	}
	uu, err := m.recursiveFindRoleAccess(ctx, page+1, limit)
	if err != nil {
		return ur, err
	}
	return append(ur, uu...), nil
}

func (m *manager) recursiveFindUserAccess(ctx context.Context, page, limit int) ([]*users.UserAccess, error) {
	ua, uat, uaerr := m.ua.FindAllUserAccess(ctx, 1, 10)
	if uaerr != nil {
		return nil, uaerr
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
