package distribution

import "time"

type Link struct {
	Id                   int64      `db:"id"`
	Hash                 string     `db:"hash"`
	PIN                  string     `db:"pin,omitempty"`
	DistributionObjectId int64      `db:"distribution_object_id"`
	Disabled             bool       `db:"disabled,omitempty"`
	Archived             bool       `db:",omitempty"`
	Published            bool       `db:"published,omitempty"`
	UsageLimit           int        `db:"usage_limit,omitempty"`
	CreatedAt            *time.Time `db:"created_at"`
	UpdatedAt            *time.Time `db:"updated_at,omitempty"`
	DisabledAt           *time.Time `db:"disabled_at,omitempty"`
	ArchivedAt           *time.Time `db:"archived_at,omitempty"`
	PublishedAt          *time.Time `db:"published_at,omitempty"`
}

type LinkVisit struct {
	Id     int64                  `db:"id"`
	LinkId int64                  `db:"link_id"`
	At     *time.Time             `db:"at"`
	Agent  string                 `db:"agent"`
	Ref    map[string]interface{} `db:"ref"`
}
