package hapi

import (
	"github.com/evorts/feednomity/domain/distribution"
	"github.com/evorts/feednomity/pkg/utils"
	"time"
)

type Distribution struct {
	Id                int64      `json:"id"`
	Topic             string     `json:"topic"`
	Disabled          bool       `json:"disabled"`
	Archived          bool       `json:"archived"`
	Distributed       bool       `json:"distributed"`
	DistributionLimit int        `json:"distribution_limit"`
	DistributionCount int        `json:"distribution_count"`
	RangeStart        *time.Time `json:"range_start"`
	RangeEnd          *time.Time `json:"range_end"`
	CreatedBy         int64      `json:"created_by"`
	UpdatedBy         int64      `json:"updated_by"`
	ForGroupId        int64      `json:"for_group_id"`
	CreatedAt         *time.Time `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at"`
	DisabledAt        *time.Time `json:"disabled_at"`
	ArchivedAt        *time.Time `json:"archived_at"`
	DistributedAt     *time.Time `json:"distributed_at"`
}

func transformDistribution(createdBy int64, items []*Distribution, excludeFields []string) (rs []*distribution.Distribution) {
	rs = make([]*distribution.Distribution, 0)
	for _, fv := range items {
		item := &distribution.Distribution{
			CreatedBy: createdBy,
		}
		if err := utils.TransformStructWithExcludes(item, fv, excludeFields); err == nil {
			rs = append(rs, item)
		}
	}
	return
}

type DistributionObject struct {
	Id               int64                         `json:"id"`
	DistributionId   int64                         `json:"distribution_id"`
	RecipientId      int64                         `json:"recipient_id"`
	RespondentId     int64                         `json:"respondent_id"`
	LinkId           int64                         `json:"link_id"`
	PublishingStatus distribution.PublishingStatus `json:"publishing_status"`
}

func transformDistributionObjects(createdBy int64, items []*DistributionObject, linksId []int64, excludeFields []string) (rs []*distribution.Object) {
	rs = make([]*distribution.Object, 0)
	for fi, fv := range items {
		item := &distribution.Object{
			CreatedBy: createdBy,
		}
		if err := utils.TransformStructWithExcludes(item, fv, excludeFields); err == nil {
			if fi+1 <= len(linksId) {
				item.LinkId = linksId[fi]
			}
			rs = append(rs, item)
		}
	}
	return
}
