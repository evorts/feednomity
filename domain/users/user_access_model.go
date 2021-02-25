package users

type AccessLevel string

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
