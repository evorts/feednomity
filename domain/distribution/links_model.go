package distribution

import "time"

type Link struct {
	Id                   int64      `json:"id"`
	Hash                 string     `json:"hash"`
	PIN                  string     `json:"pin,omitempty"`
	DistributionObjectId int64      `json:"distribution_object_id"`
	Disabled             bool       `json:"disabled,omitempty"`
	Archived             bool       `json:",omitempty"`
	Published            bool       `json:"published,omitempty"`
	UsageLimit           int        `json:"usage_limit,omitempty"`
	CreatedAt            *time.Time `json:"created_at"`
	UpdatedAt            *time.Time `json:"updated_at,omitempty"`
	DisabledAt           *time.Time `json:"disabled_at,omitempty"`
	ArchivedAt           *time.Time `json:"archived_at,omitempty"`
	PublishedAt          *time.Time `json:"published_at,omitempty"`
}

type LinkVisit struct {
	Id     int64                  `json:"id"`
	LinkId int64                  `json:"link_id"`
	At     *time.Time             `json:"at"`
	Agent  string                 `json:"agent"`
	Ref    map[string]interface{} `json:"ref"`
}
