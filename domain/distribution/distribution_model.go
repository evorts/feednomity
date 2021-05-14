package distribution

import "time"

type Distribution struct {
	Id                int64      `db:"id"`
	Topic             string     `db:"topic"`
	Disabled          bool       `db:"disabled"`
	Archived          bool       `db:"archived"`
	Distributed       bool       `db:"distributed"`
	DistributionLimit int        `db:"distribution_limit"`
	DistributionCount int        `db:"distribution_count"`
	RangeStart        *time.Time `db:"range_start"`
	RangeEnd          *time.Time `db:"range_end"`
	CreatedBy         int64      `db:"created_by"`
	ForGroupId        int64      `db:"for_group_id"`
	CreatedAt         *time.Time `db:"created_at"`
	UpdatedAt         *time.Time `db:"updated_at"`
	DisabledAt        *time.Time `db:"disabled_at"`
	ArchivedAt        *time.Time `db:"archived_at"`
	DistributedAt     *time.Time `db:"distributed_at"`
}

type PublishingStatus string

const (
	PublishingNone   PublishingStatus = "none"
	PublishingSent   PublishingStatus = "sent"
	PublishingFailed PublishingStatus = "failed"
)

type Object struct {
	Id               int64                    `db:"id"`
	DistributionId   int64                    `db:"distribution_id"`
	RecipientId      int64                    `db:"recipient_id"`
	RespondentId     int64                    `db:"respondent_ids"`
	LinkId           int64                    `db:"link_id"`
	PublishingStatus PublishingStatus         `db:"publishing_status"`
	PublishingLog    []map[string]interface{} `db:"publishing_log"`
	RetryCount       int                      `db:"retry_count"`
	CreatedBy        int64                    `db:"created_by"`
	UpdatedBy        int64                    `db:"updated_by"`
	CreatedAt        *time.Time               `db:"created_at"`
	UpdatedAt        *time.Time               `db:"updated_at"`
	PublishedAt      *time.Time               `db:"published_at"`
}

type Queue struct {
	Id                   int64                  `db:"id"`
	DistributionObjectId int64                  `db:"distribution_object_id"`
	RecipientId          int64                  `db:"recipient_id"`
	RespondentId         int64                  `db:"respondent_id"`
	FromEmail            string                 `db:"from_email"`
	ToEmail              string                 `db:"to_email"`
	Subject              string                 `db:"subject"`
	Template             string                 `db:"template"`
	Arguments            map[string]interface{} `db:"arguments"`
}

type Log struct {
	Id         int64                  `db:"id"`
	Action     string                 `db:"action"`
	Values     map[string]interface{} `db:"values"`
	ValuesPrev map[string]interface{} `db:"values_prev"`
	Notes      string                 `db:"notes"`
	At         *time.Time             `db:"at"`
}
