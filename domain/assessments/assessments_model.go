package assessments

import (
	"github.com/evorts/feednomity/pkg/utils"
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

func (f *Factor) Update(key string, rating int, note string) bool {
	if f.Key == key {
		f.Rating = rating
		f.Note = note
		return true
	}
	if len(f.Items) < 1 {
		return false
	}
	for _, item := range f.Items {
		if item.Update(key, rating, note) {
			return true
		}
	}
	return false
}

func BindToFeedbackFactors(parentKey string, value map[string]interface{}, factors *Factor) {
	if len(parentKey) > 0 {
		vr, ok := value["rating"]
		if ok {
			rating, ok2 := vr.(float64)
			if !ok2 {
				return
			}
			factors.Update(parentKey, int(rating), utils.IIf(value["note"] == nil, "", value["note"].(string)))
			return
		}
	}
	for k, v := range value {
		vv, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		BindToFeedbackFactors(k, vv, factors)
	}
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
		Values    []int      `yaml:"values" db:"values"`
		Labels    []string   `yaml:"labels" db:"labels"`
		Threshold [][]string `yaml:"threshold" db:"threshold"`
	} `yaml:"ratings" db:"ratings"`
	StrengthsFieldCount    int     `yaml:"strengths_field_count" db:"strengths_field_count"`
	ImprovementsFieldCount int     `yaml:"improvements_field_count" db:"improvements_field_count"`
	Factors                *Factor `yaml:"factors" db:"factors"`
}
