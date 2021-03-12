package users

import "time"

type User struct {
	Id          int64      `json:"id"`
	Username    string     `json:"username"`
	DisplayName string     `json:"display_name"`
	Email       string     `json:"email"`
	Phone       string     `json:"phone"`
	Password    string     `json:"password"`
	Role        UserRole   `json:"role"`
	GroupId     int        `json:"group_id"`
	CreatedDate *time.Time `json:"created_date"`
	UpdatedDate *time.Time `json:"updated_date"`
}

type Group struct {
	Id         int64      `json:"id"`
	Name       string     `json:"name"`
	Disabled   bool       `json:"disabled"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	DisabledAt *time.Time `json:"disabled_at"`
}
