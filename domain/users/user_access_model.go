package users

import (
	"database/sql/driver"
	"time"
)

type AccessMethod string

func (a *AccessMethod) String() string {
	if a == nil {
		return ""
	}
	return string(*a)
}
func (a *AccessMethod) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return a.String(), nil
}

func (a *AccessMethod) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	v, ok := src.(string)
	if ok {
		*a = AccessMethod(v)
	}
	return nil
}

func AccessMethodsToStringArray(al []AccessMethod) []string {
	rs := make([]string, 0)
	for _, v := range al {
		rs = append(rs, v.String())
	}
	return rs
}

const (
	AccessMethodRead     AccessMethod = "get"
	AccessMethodInsert   AccessMethod = "post"
	AccessMethodUpdate   AccessMethod = "put"
	AccessMethodDelete   AccessMethod = "delete"
	AccessMethodHeadInfo AccessMethod = "head"
	AccessMethodOptions  AccessMethod = "options"
)

type UserRole string

const (
	UserRoleSysAdmin   UserRole = "sysadmin"   // master super user, expected to be only one
	UserRoleSiteAdmin  UserRole = "site-admin" // similar to sysadmin, but can't remove sysadmin
	UserRoleAdmin      UserRole = "admin"
	UserRoleSupervisor UserRole = "supervisor"
	UserRoleMember     UserRole = "member"
	UserRoleInvitation UserRole = "guest"
	UserRoleCustom     UserRole = "custom"
)

func (u UserRole) String() string {
	return string(u)
}

type UserRoleAccess struct {
	Id               int64          `db:"id"`
	Role             UserRole       `db:"role"`
	Path             string         `db:"path"`
	Regex            bool           `db:"regex"`
	Disabled         bool           `db:"disabled"`
	AccessAllowed    []AccessMethod `db:"access_allowed"`
	AccessDisallowed []AccessMethod `db:"access_disallowed"`
	AccessScope      string         `db:"access_scope"`
	CreatedAt        *time.Time     `db:"created_at"`
	UpdatedAt        *time.Time     `db:"updated_at"`
	DisabledAt       *time.Time     `db:"disabled_at"`
}

type UserAccess struct {
	Id               int64          `db:"id"`
	UserId           int64          `db:"user_id"`
	Path             string         `db:"path"`
	Regex            bool           `db:"regex"`
	Disabled         bool           `db:"disabled"`
	AccessAllowed    []AccessMethod `db:"access_allowed"`
	AccessDisallowed []AccessMethod `db:"access_disallowed"`
	AccessScope      string         `db:"access_scope"`
	CreatedAt        *time.Time     `db:"created_at"`
	UpdatedAt        *time.Time     `db:"updated_at"`
	DisabledAt       *time.Time     `db:"disabled_at"`
}
