package users

import (
	"regexp"
	"time"
)

var (
	validPINPattern      = regexp.MustCompile("\\d{6}")
	validPasswordPattern = regexp.MustCompile("\\w{8}")
)

type PIN string

func (p PIN) Rule() string {
	return "must be numeric with length of 6"
}

func (p PIN) Value() string {
	return string(p)
}

func (p PIN) Valid() bool {
	return validPINPattern.MatchString(p.Value())
}

type PASSWORD string

func (p PASSWORD) Rule() string {
	return "must be alpha-numeric and at least 8 character"
}

func (p PASSWORD) Value() string {
	return string(p)
}

func (p PASSWORD) Valid() bool {
	return validPasswordPattern.MatchString(p.Value())
}

type User struct {
	Id          int64                  `db:"id"`
	Username    string                 `db:"username"`
	DisplayName string                 `db:"display_name"`
	Attributes  map[string]interface{} `db:"attributes"`
	Email       string                 `db:"email"`
	Phone       string                 `db:"phone"`
	Password    string                 `db:"password"`
	PIN         string                 `db:"pin"`
	AccessRole  UserRole               `db:"access_role"`
	JobRole     string                 `db:"job_role"`
	Assignment  string                 `db:"assignment"`
	GroupId     int64                  `db:"group_id"`
	Disabled    bool                   `db:"disabled"`
	CreatedAt   *time.Time             `db:"created_at"`
	UpdatedAt   *time.Time             `db:"updated_at"`
	DisabledAt  *time.Time             `db:"disabled_at"`
}

type Group struct {
	Id         int64      `db:"id"`
	Name       string     `db:"name"`
	OrgId      int64        `db:"org_id"`
	Disabled   bool       `db:"disabled"`
	CreatedAt  *time.Time `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
	DisabledAt *time.Time `db:"disabled_at"`
}

type Organization struct {
	Id         int64        `db:"id"`
	Name       string     `db:"name"`
	Address    string     `db:"address"`
	Phone      string     `db:"phone"`
	Disabled   bool       `db:"disabled"`
	CreatedAt  *time.Time `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
	DisabledAt *time.Time `db:"disabled_at"`
}
