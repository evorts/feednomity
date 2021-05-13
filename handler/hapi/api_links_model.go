package hapi

import (
	"github.com/evorts/feednomity/domain/distribution"
	"github.com/evorts/feednomity/pkg/utils"
	"time"
)

type Link struct {
	Id          int64                  `json:"id"`
	Hash        string                 `json:"hash"`
	PIN         string                 `json:"pin,omitempty"`
	Disabled    bool                   `json:"disabled,omitempty"`
	Archived    bool                   `json:",omitempty"`
	Published   bool                   `json:"published,omitempty"`
	UsageLimit  int                    `json:"usage_limit,omitempty"`
	Attributes  map[string]interface{} `json:"attributes"`
	CreatedAt   *time.Time             `json:"created_at"`
	UpdatedAt   *time.Time             `json:"updated_at,omitempty"`
	DisabledAt  *time.Time             `json:"disabled_at,omitempty"`
	ArchivedAt  *time.Time             `json:"archived_at,omitempty"`
	PublishedAt *time.Time             `json:"published_at,omitempty"`
}

func transformLinks(items []*Link, expireAt *time.Time, excludeFields []string) (rs []*distribution.Link) {
	rs = make([]*distribution.Link, 0)
	for _, fv := range items {
		item := &distribution.Link{
			ExpiredAt: expireAt,
		}
		if err := utils.TransformStructWithExcludes(item, fv, excludeFields); err == nil {
			rs = append(rs, item)
		}
	}
	return
}

type HashData struct {
	ExpireAt   time.Time   `json:"expire_at"`
	RealHash   string      `json:"real_hash"`
	Attributes interface{} `json:"attributes"`
}
