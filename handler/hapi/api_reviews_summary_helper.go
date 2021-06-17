package hapi

import (
	"encoding/json"
	"fmt"
	"github.com/evorts/feednomity/domain/assessments"
	"github.com/evorts/feednomity/domain/feedbacks"
	"github.com/evorts/feednomity/handler/helpers"
	"github.com/evorts/feednomity/pkg/utils"
	"sort"
)

func generateReviewSummaryData(feeds []*feedbacks.Feedback, factors *assessments.Template, excludeStatus ...feedbacks.Status) (items []*FeedbackSummaryResponseItem, err error) {
	items = make([]*FeedbackSummaryResponseItem, 0)
	var (
		mapByRecipient = make(map[int64]*FeedbackSummaryResponseItem, 0) // map by recipient
		eval           = utils.NewEval()
	)
	loopFeed:
	for _, feed := range feeds {
		if excludeStatus != nil {
			for _, wStatus := range excludeStatus {
				if wStatus == feed.Status {
					continue loopFeed
				}
			}
		}
		if feed.Status != feedbacks.StatusFinal {
			fmt.Println(feed.Id, feed.RecipientUsername, feed.RespondentUsername, feed.Status)
		}
		cnt, okc := feed.Content["raw"]
		if !okc || cnt == nil {
			continue
		}
		cntB, err2 := json.Marshal(cnt)
		if err2 != nil {
			continue
		}
		var content map[string]interface{}
		if err = json.Unmarshal(cntB, &content); err != nil {
			continue
		}
		var factor assessments.Factor
		assessments.BindToFeedbackFactors("", content, factors.Factors)
		fByte, _ := json.Marshal(factors.Factors)
		_ = json.Unmarshal(fByte, &factor)
		score := helpers.CalculateScore(&factor)
		item, ok := mapByRecipient[feed.RecipientId]
		if !ok {
			mapByRecipient[feed.RecipientId] = &FeedbackSummaryResponseItem{
				DistributionId: feed.DistributionId,
				Recipient: Person{
					Id:         feed.RecipientId,
					Name:       feed.RecipientName,
					GroupId:    feed.RecipientGroupId,
					GroupName:  feed.RecipientGroupName,
					OrgId:      feed.RecipientOrgId,
					OrgName:    feed.RecipientOrgName,
					Role:       feed.RecipientRole,
					Assignment: feed.RecipientAssignment,
				},
				TotalScore: 0,
				Rating:     "",
				RangeStart: feed.RangeStart,
				RangeEnd:   feed.RangeEnd,
				Items:      make([]*FeedbackItem, 0),
			}
			item = mapByRecipient[feed.RecipientId]
		}
		item.TotalScore = (item.TotalScore + score) / utils.IIfF64(item.TotalScore == 0, 1, 2)
		item.Rating = helpers.GetRating(eval, factors.Ratings.Labels, factors.Ratings.Threshold, item.TotalScore)
		item.Items = append(item.Items, &FeedbackItem{
			Id: feed.Id,
			Respondent: Person{
				Id:         feed.RespondentId,
				Name:       feed.RespondentName,
				GroupId:    feed.RespondentGroupId,
				GroupName:  feed.RespondentGroupName,
				OrgId:      feed.RespondentOrgId,
				OrgName:    feed.RespondentOrgName,
				Role:       feed.RespondentRole,
				Assignment: feed.RespondentAssignment,
			},
			DistributionObjectId: feed.DistributionObjectId,
			Score:                score,
			Rating:               helpers.GetRating(eval, factors.Ratings.Labels, factors.Ratings.Threshold, score),
			Factors:              &factor,
			FactorsRaw:           content,
			Strengths:            utils.ArrayInterface(content["strengths"].([]interface{})).ToArrayString().Reduce(),
			NeedImprovements:     utils.ArrayInterface(content["improves"].([]interface{})).ToArrayString().Reduce(),
			Status:               feed.Status,
			UpdatedAt:            feed.UpdatedAt,
		})
	}
	//map to array
	for _, item := range mapByRecipient {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Recipient.Name < items[j].Recipient.Name
	})
	return
}
