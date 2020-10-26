package users

import "time"

type UserRole string

const (
	UserRoleSysAdmin   UserRole = "sysadmin"
	UserRoleAdmin      UserRole = "admin"
	UserRoleMember     UserRole = "member"
	UserRoleInvitation UserRole = "invitation"
	UserRoleCustom     UserRole = "custom"
)

type RequestMethod string

const (
	RequestMethodGet     RequestMethod = "get"
	RequestMethodPost    RequestMethod = "post"
	RequestMethodPut     RequestMethod = "put"
	RequestMethodDelete  RequestMethod = "delete"
	RequestMethodHead    RequestMethod = "head"
	RequestMethodOptions RequestMethod = "options"
)

type AccessLevel string

const (
	AccessLevelReadOnly  AccessLevel = "ro"
	AccessLevelWriteOnly AccessLevel = "wo"
	AccessLevelReadWrite AccessLevel = "rw"
)

type AccessScope string

const (
	AccessScopeCustom AccessScope = "custom"
	AccessScopeAll    AccessScope = "all"
)

type User struct {
	Id          int64
	Username    string
	DisplayName string
	Email       string
	Phone       string
	Password    string
	Role        UserRole
	CreatedDate time.Time
	UpdatedDate time.Time
}

type UserRoleAccess struct {
	Id            int64
	Role          UserRole
	Path          string
	MethodAllowed []RequestMethod
	Disabled      bool
	AccessLevel   AccessLevel
}

type UserAccess struct {
	Id            int64
	UserId        int64
	Scope         AccessScope
	Path          string
	MethodAllowed []RequestMethod
	Disabled      bool
	AccessLevel   AccessLevel
}
