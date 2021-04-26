package assessments

import (
	"html/template"
	"time"
)

type Factor struct {
	Key         string        `db:"key"`
	Title       string        `db:"title"`
	Description template.HTML `db:"description"`
	Weight      float32       `db:"weight"`
	Rating      int           `db:"rating"`
	Note        string        `db:"note"`
	Items       []*Factor     `db:"items"`
}

type Client struct {
	Id           int                    `db:"id"`
	Name         string                 `db:"name"`
	Email        string                 `db:"email"`
	Phone        string                 `db:"phone"`
	Role         string                 `db:"role"`
	Assignment   string                 `db:"assignment"`
	Organization string                 `db:"organization"`
	Group        string                 `db:"group"`
	GroupId      string                 `db:"group_id"`
	Attributes   map[string]interface{} `db:"attributes"`
}

type Item struct {
	Recipient        Client     `db:"recipient"`
	Respondent       Client     `db:"respondent"`
	PeriodSince      *time.Time `db:"period_since"`
	PeriodUntil      *time.Time `db:"period_until"`
	Ratings          []int      `db:"ratings"`
	RatingsLabel     []string   `db:"ratings_label"`
	Factors          *Factor    `db:"assessment_factors"`
	Strengths        []string   `db:"strengths"`
	NeedImprovements []string   `db:"need_improvements"`
}

type Template struct {
	Ratings struct {
		Values []int    `yaml:"values" db:"values"`
		Labels []string `yaml:"labels" db:"labels"`
	} `yaml:"ratings" db:"ratings"`
	StrengthsFieldCount    int     `yaml:"strengths_field_count" db:"strengths_field_count"`
	ImprovementsFieldCount int     `yaml:"improvements_field_count" db:"improvements_field_count"`
	Factors                *Factor `yaml:"factors" db:"factors"`
}
