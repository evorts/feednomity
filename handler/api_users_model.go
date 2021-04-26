package handler

import (
	"github.com/evorts/feednomity/domain/users"
	"github.com/evorts/feednomity/pkg/utils"
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
	GroupId     int                    `json:"group_id"`
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
