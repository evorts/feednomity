package users

import "database/sql/driver"

type AccessLevel string

func (a *AccessLevel) String() string {
	if a == nil {
		return ""
	}
	return string(*a)
}
func (a *AccessLevel) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return a.String(), nil
}

func (a *AccessLevel) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	v, ok := src.(string)
	if ok {
		*a = AccessLevel(v)
	}
	return nil
}

const (
	AccessLevelRead     AccessLevel = "get"
	AccessLevelInsert   AccessLevel = "post"
	AccessLevelUpdate   AccessLevel = "put"
	AccessLevelDelete   AccessLevel = "delete"
	AccessLevelHeadInfo AccessLevel = "head"
	AccessLevelOptions  AccessLevel = "options"
)

type UserRole string

const (
	UserRoleSysAdmin   UserRole = "sysadmin" // master super user, expected to be only one
	UserRoleSiteAdmin  UserRole = "site-admin" // similar to sysadmin, but can't remove sysadmin
	UserRoleAdmin      UserRole = "admin"
	UserRoleSupervisor UserRole = "supervisor"
	UserRoleMember     UserRole = "member"
	UserRoleInvitation UserRole = "guest"
	UserRoleCustom     UserRole = "custom"
)

type UserRoleAccess struct {
	Id            int64
	Role          UserRole
	Path          string
	Disabled      bool
	AccessAllowed []AccessLevel
}

type UserAccess struct {
	Id               int64
	UserId           int64
	Path             string
	Disabled         bool
	AccessAllowed    []AccessLevel
	AccessDisallowed []AccessLevel
}
