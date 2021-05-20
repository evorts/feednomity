package hapi

import (
	"github.com/evorts/feednomity/domain/users"
	"github.com/evorts/feednomity/pkg/utils"
	"time"
)

type UserGroupRequest struct {
	Id         int64      `json:"id"`
	Name       string     `json:"name"`
	OrgId      int64      `json:"org_id"`
	Disabled   bool       `json:"disabled"`
}

type UserGroupResponse struct {
	Id         int64      `json:"id"`
	Name       string     `json:"name"`
	OrgId      int64      `json:"org_id"`
	Disabled   bool       `json:"disabled"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	DisabledAt *time.Time `json:"disabled_at"`
}

func transformGroups(f []*UserGroupRequest) (t []*users.Group) {
	t = make([]*users.Group, 0)
	for _, fv := range f {
		u := &users.Group{}
		if err := utils.TransformStruct(u, fv); err == nil {
			t = append(t, u)
		}
	}
	return
}

func transformGroupsReverse(f []*users.Group) (t []*UserGroupResponse) {
	t = make([]*UserGroupResponse, 0)
	for _, fv := range f {
		u := &UserGroupResponse{}
		if err := utils.TransformStruct(u, fv); err == nil {
			t = append(t, u)
		}
	}
	return
}