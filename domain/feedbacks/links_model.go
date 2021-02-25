package feedbacks

import "time"

type Link struct {
	Id          int64      `json:"id"`
	Hash        Hash       `json:"hash"`
	PIN         PIN     `json:"pin"`
	GroupId     int64      `json:"group_id"`
	Disabled    bool       `json:"disabled"`
	Published   bool       `json:"published"`
	UsageLimit  int64      `json:"usage_limit"`
	CreatedAt   *time.Time `json:"-"`
	UpdatedAt   *time.Time `json:"-"`
	DisabledAt  *time.Time `json:"-"`
	PublishedAt *time.Time `json:"-"`
}

type LinkVisit struct {
	Id     int64
	LinkId int64
	At     *time.Time
	Agent  string
	Ref    map[string]interface{}
}