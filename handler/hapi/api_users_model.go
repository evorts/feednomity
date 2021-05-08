package hapi

import (
	"github.com/evorts/feednomity/domain/users"
	"github.com/evorts/feednomity/pkg/utils"
	"time"
)

type User struct {
	Id          int64                  `json:"id"`
	Username    string                 `json:"username"`
	DisplayName string                 `json:"display_name"`
	Attributes  map[string]interface{} `json:"attributes"`
	Email       string                 `json:"email"`
	Phone       string                 `json:"phone"`
	Password    string                 `json:"password"`
	PIN         string                 `json:"pin"`
	AccessRole  users.UserRole         `json:"access_role"`
	JobRole     string                 `json:"job_role"`
	Assignment  string                 `json:"assignment"`
	GroupId     int64                  `json:"group_id"`
	Disabled    bool                   `json:"disabled"`
}

func transformUsers(f []*User) (t []*users.User) {
	t = make([]*users.User, 0)
	for _, fv := range f {
		u := &users.User{}
		if err := utils.TransformStruct(u, fv); err == nil {
			t = append(t, u)
		}
	}
	return
}

type UserGroup struct {
	Id         int64      `json:"id"`
	Name       string     `json:"name"`
	OrgId      int64      `json:"org_id"`
	Disabled   bool       `json:"disabled"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	DisabledAt *time.Time `json:"disabled_at"`
}

func transformGroups(f []*UserGroup) (t []*users.Group) {
	t = make([]*users.Group, 0)
	for _, fv := range f {
		u := &users.Group{}
		if err := utils.TransformStruct(u, fv); err == nil {
			t = append(t, u)
		}
	}
	return
}

type UserOrg struct {
	Id         int64      `json:"id"`
	Name       string     `json:"name"`
	Address    string     `json:"address"`
	Phone      string     `json:"phone"`
	Disabled   bool       `json:"disabled"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	DisabledAt *time.Time `json:"disabled_at"`
}

func transformOrganizations(f []*UserOrg) (t []*users.Organization) {
	t = make([]*users.Organization, 0)
	for _, fv := range f {
		u := &users.Organization{}
		if err := utils.TransformStruct(u, fv); err == nil {
			t = append(t, u)
		}
	}
	return
}