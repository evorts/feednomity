package hapi

import (
	"github.com/evorts/feednomity/domain/users"
	"github.com/evorts/feednomity/pkg/utils"
	"time"
)

type UserRequest struct {
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

type UserResponse struct {
	Id             int64                  `json:"id"`
	Username       string                 `json:"username"`
	DisplayName    string                 `json:"display_name"`
	Attributes     map[string]interface{} `json:"attributes"`
	Email          string                 `json:"email"`
	Phone          string                 `json:"phone"`
	AccessRole     users.UserRole         `json:"access_role"`
	JobRole        string                 `json:"job_role"`
	Assignment     string                 `json:"assignment"`
	GroupId        int64                  `json:"group_id"`
	OrganizationId int64                  `json:"-"`
	Disabled       bool                   `json:"disabled"`
	CreatedAt      *time.Time             `json:"created_at"`
	UpdatedAt      *time.Time             `json:"updated_at"`
	DisabledAt     *time.Time             `json:"disabled_at"`
}

func transformUsers(f []*UserRequest) (t []*users.User) {
	t = make([]*users.User, 0)
	for _, fv := range f {
		u := &users.User{}
		if err := utils.TransformStruct(u, fv); err == nil {
			t = append(t, u)
		}
	}
	return
}

func transformUsersReverse(f []*users.User) (t []*UserResponse) {
	t = make([]*UserResponse, 0)
	for _, fv := range f {
		u := &UserResponse{}
		if err := utils.TransformStruct(u, fv); err == nil {
			t = append(t, u)
		}
	}
	return
}
