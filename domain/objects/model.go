package objects

import (
	"time"
)

// Object /** storing persons or employees in this struct **/
type Object struct {
	Id          int64                  `db:"id"`
	Name        string                 `db:"name"`
	Attributes  map[string]interface{} `db:"attributes"`
	Email       string                 `db:"email"`
	Phone       string                 `db:"phone"`
	Role        string                 `db:"role"`
	Assignment  string                 `db:"assignment"`
	UserGroupId int64                  `db:"user_group_id"` // the recipient record belong to which user group
	Disabled    bool                   `db:"disabled"`
	Archived    bool                   `db:"archived"`
	CreatedAt   *time.Time             `db:"created_at"`
	UpdatedAt   *time.Time             `db:"updated_at"`
	DisabledAt  *time.Time             `db:"disabled_at"`
	ArchivedAt  *time.Time             `db:"archived_at"`
}
