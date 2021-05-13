package distribution

import "time"

type Link struct {
	Id          int64                  `db:"id"`
	Hash        string                 `db:"hash"`
	PIN         string                 `db:"pin,omitempty"`
	Disabled    bool                   `db:"disabled,omitempty"`
	Archived    bool                   `db:",omitempty"`
	Published   bool                   `db:"published,omitempty"`
	UsageLimit  int                    `db:"usage_limit,omitempty"`
	Attributes  map[string]interface{} `db:"attributes"`
	CreatedBy   int64                  `db:"created_by"`
	UpdatedBy   int64                  `db:"updated_by,omitempty"`
	ExpiredAt   *time.Time             `db:"expired_at"`
	CreatedAt   *time.Time             `db:"created_at"`
	UpdatedAt   *time.Time             `db:"updated_at,omitempty"`
	DisabledAt  *time.Time             `db:"disabled_at,omitempty"`
	ArchivedAt  *time.Time             `db:"archived_at,omitempty"`
	PublishedAt *time.Time             `db:"published_at,omitempty"`
}

type LinkVisit struct {
	Id     int64                  `db:"id"`
	LinkId int64                  `db:"link_id"`
	By     int64                  `db:"by"`
	ByName string                 `db:"by_name"`
	At     *time.Time             `db:"at"`
	Agent  string                 `db:"agent"`
	Ref    map[string]interface{} `db:"ref"`
}
