package feedbacks

import "time"

type Audience struct {
	Id         int64
	Title      string
	Emails     []string
	Disabled   bool
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	DisabledAt *time.Time
}

