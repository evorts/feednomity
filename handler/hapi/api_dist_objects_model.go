package hapi

import (
	"github.com/evorts/feednomity/domain/distribution"
	"github.com/evorts/feednomity/pkg/utils"
	"time"
)

type DistributionObjectRequest struct {
	Id               int64                         `json:"id"`
	DistributionId   int64                         `json:"distribution_id"`
	RecipientId      int64                         `json:"recipient_id"`
	RespondentId     int64                         `json:"respondent_id"`
	LinkId           int64                         `json:"link_id"`
	PublishingStatus distribution.PublishingStatus `json:"publishing_status"`
}

type DistributionObjectResponse struct {
	Id               int64                         `json:"id"`
	DistributionId   int64                         `json:"distribution_id"`
	RecipientId      int64                         `json:"recipient_id"`
	RespondentId     int64                         `json:"respondent_ids"`
	LinkId           int64                         `json:"link_id"`
	PublishingStatus distribution.PublishingStatus `json:"publishing_status"`
	PublishingLog    []map[string]interface{}      `json:"publishing_log"`
	RetryCount       int64                         `json:"retry_count"`
	CreatedBy        int64                         `json:"created_by"`
	UpdatedBy        int64                         `json:"updated_by"`
	CreatedAt        *time.Time                    `json:"created_at"`
	UpdatedAt        *time.Time                    `json:"updated_at"`
	PublishedAt      *time.Time                    `json:"published_at"`
}

func transformDistributionObjects(createdBy int64, items []*DistributionObjectRequest, linksId []int64, excludeFields []string) (rs []*distribution.Object) {
	rs = make([]*distribution.Object, 0)
	for fi, fv := range items {
		item := &distribution.Object{
			CreatedBy: createdBy,
		}
		if err := utils.TransformStructWithExcludes(item, fv, excludeFields, false); err == nil {
			if fi+1 <= len(linksId) {
				item.LinkId = linksId[fi]
			}
			rs = append(rs, item)
		}
	}
	return
}

func transformDistributionObjectsReverse(items []*distribution.Object) (rs []*DistributionObjectResponse) {
	rs = make([]*DistributionObjectResponse, 0)
	for _, fv := range items {
		item := &DistributionObjectResponse{}
		if err := utils.TransformStruct(item, fv); err == nil {
			rs = append(rs, item)
		}
	}
	return
}
