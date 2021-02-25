package objects

import (
	"time"
)

/** storing persons or employees in this struct **/
type Object struct {
	Id          int                    `json:"id"`
	Name        string                 `json:"name"`
	Attributes  map[string]interface{} `json:"attributes"`
	Email       string                 `json:"email"`
	Phone       string                 `json:"phone"`
	Role        string                 `json:"role"`
	Assignment  string                 `json:"assignment"`
	UserGroupId int                    `json:"user_group_id"` // the recipient record belong to which user group
	Disabled    bool                   `json:"disabled"`
	Archived    bool                   `json:"archived"`
	CreatedAt   *time.Time             `json:"created_at"`
	UpdatedAt   *time.Time             `json:"updated_at"`
	DisabledAt  *time.Time             `json:"disabled_at"`
	ArchivedAt  *time.Time             `json:"archived_at"`
}