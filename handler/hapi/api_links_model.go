package hapi

import (
	"github.com/evorts/feednomity/domain/distribution"
	"github.com/evorts/feednomity/pkg/utils"
	"github.com/segmentio/ksuid"
	"time"
)

type Link struct {
	Id         int64                  `json:"id"`
	PIN        string                 `json:"pin,omitempty"`
	Disabled   bool                   `json:"disabled,omitempty"`
	Archived   bool                   `json:"archived,omitempty"`
	Published  bool                   `json:"published,omitempty"`
	UsageLimit int                    `json:"usage_limit,omitempty"`
	Attributes map[string]interface{} `json:"attributes"`
	ExpireAt   *time.Time             `json:"expire_at"`
}

func transformLinks(items []*Link, expireAt *time.Time, excludeFields []string) (rs []*distribution.Link) {
	rs = make([]*distribution.Link, 0)
	for _, fv := range items {
		item := &distribution.Link{
			ExpiredAt: expireAt,
			Hash:      ksuid.New().String(),
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
