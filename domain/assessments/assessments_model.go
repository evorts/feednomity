package assessments

import (
	"html/template"
	"time"
)

type Rating int

const (
	RatingLowest  Rating = 1
	RatingLow     Rating = 2
	RatingMedium  Rating = 3
	RatingHigh    Rating = 4
	RatingHighest Rating = 5
)

type Factor struct {
	Key         string        `json:"key"`
	Title       string        `json:"title"`
	Description template.HTML `json:"description"`
	Weight      float32       `json:"weight"`
	Rating      Rating        `json:"rating"`
	Notes       string        `json:"notes"`
	Items       []Factor      `json:"items"`
}

type Client struct {
	Id           int                    `json:"id"`
	Name         string                 `json:"name"`
	Email        string                 `json:"email"`
	Phone        string                 `json:"phone"`
	Role         string                 `json:"role"`
	Assignment   string                 `json:"assignment"`
	Organization string                 `json:"organization"`
	Group        string                 `json:"group"`
	GroupId      string                 `json:"group_id"`
	Attributes   map[string]interface{} `json:"attributes"`
}

type Item struct {
	Recipient        Client     `json:"recipient"`
	Respondent       Client     `json:"respondent"`
	PeriodSince      *time.Time `json:"period_since"`
	PeriodUntil      *time.Time `json:"period_until"`
	Ratings          []int      `json:"ratings"`
	RatingsLabel     []string   `json:"ratings_label"`
	Factors          *Factor    `json:"assessment_factors"`
	Strengths        []string   `json:"strengths"`
	NeedImprovements []string   `json:"need_improvements"`
}

type Template struct {
	Ratings struct {
		Values []int    `yaml:"values" json:"values"`
		Labels []string `yaml:"labels" json:"labels"`
	} `yaml:"ratings" json:"ratings"`
	Factors Factor `yaml:"factors" json:"factors"`
}
