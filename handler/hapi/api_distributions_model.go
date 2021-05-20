package hapi

import (
	"github.com/evorts/feednomity/domain/distribution"
	"github.com/evorts/feednomity/pkg/utils"
	"time"
)

type DistributionRequest struct {
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
}

type DistributionResponse struct {
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
	ForGroupId        int64      `json:"for_group_id"`
	CreatedAt         *time.Time `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at"`
	DisabledAt        *time.Time `json:"disabled_at"`
	ArchivedAt        *time.Time `json:"archived_at"`
	DistributedAt     *time.Time `json:"distributed_at"`
}

func transformDistribution(createdBy int64, items []*DistributionRequest, excludeFields []string) (rs []*distribution.Distribution) {
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

func transformDistributionReverse(items []*distribution.Distribution) (rs []*DistributionResponse) {
	rs = make([]*DistributionResponse, 0)
	for _, fv := range items {
		item := &DistributionResponse{}
		if err := utils.TransformStruct(item, fv); err == nil {
			rs = append(rs, item)
		}
	}
	return
}
