package feedbacks

import "html/template"

type AssessmentFactor struct {
	Seq         func(i int) int
	Key         string             `json:"key"`
	Title       string             `json:"title"`
	Description template.HTML      `json:"description"`
	Weight      float32            `json:"weight"`
	Rating      int                `json:"rating"`
	Notes       string             `json:"notes"`
	Items       []AssessmentFactor `json:"items"`
}
