package distribution

import "time"

type Distribution struct {
	Id                int64      `json:"id"`
	Topic             string     `json:"topic"`
	Disabled          bool       `json:"disabled"`
	Archived          bool       `json:"archived"`
	Distributed       bool       `json:"distributed"`
	DistributionLimit int        `json:"distribution_limit"`
	DistributionCount int        `json:"distribution_count"`
	CreatedAt         *time.Time `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at"`
	DisabledAt        *time.Time `json:"disabled_at"`
	ArchivedAt        *time.Time `json:"archived_at"`
	DistributedAt     *time.Time `json:"distributed_at"`
}

type Object struct {
	Id               int64                  `json:"id"`
	DistributionId   int64                  `json:"distribution_id"`
	RecipientId      int64                  `json:"recipient_id"`
	RespondentId     int64                  `json:"respondent_ids"`
	LinkId           int64                  `json:"link_id"`
	PublishingStatus string                 `json:"publishing_status"`
	PublishingLog    []map[string]interface{} `json:"publishing_log"`
	CreatedAt        *time.Time             `json:"created_at"`
	UpdatedAt        *time.Time             `json:"updated_at"`
	PublishedAt      *time.Time             `json:"published_at"`
}

type Log struct {
	Id         int64                  `json:"id"`
	Action     string                 `json:"action"`
	Values     map[string]interface{} `json:"values"`
	ValuesPrev map[string]interface{} `json:"values_prev"`
	Notes      string                 `json:"notes"`
	At         *time.Time             `json:"at"`
}
