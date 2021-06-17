package hapi

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/evorts/feednomity/domain/assessments"
	"github.com/evorts/feednomity/domain/feedbacks"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/spreadsheets"
	"github.com/evorts/feednomity/pkg/utils"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
	"strings"
)

func ApiSummaryReviewsExport(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).PrepareRestful()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.IManager)
	ds := req.GetContext().Get("db").(database.IManager)

	log.Log("api_summary_export_handler", "request received")

	var payload struct {
		FileType       string              `json:"file_type"`
		ExcludeStatus  []feedbacks.Status `json:"exclude_status"`
		DistributionId int64               `json:"distribution_id"`
	}

	err := req.UnmarshallBody(&payload)

	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "SME:ERR:BND",
				Message: "Bad Request! Something wrong with the payload of your request.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}

	//validation
	errs := make(map[string]string, 0)

	if !utils.InArray([]interface{}{"xls", "xlsx"}, strings.ToLower(payload.FileType)) {
		errs["filetype"] = "Incorrect filetype"
	}
	if payload.DistributionId < 1 {
		errs["distribution"] = "No distribution argument supplied"
	}

	if len(errs) > 0 {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "SME:ERR:VAL",
				Message: "Bad Request! Your request resulting validation error.",
				Reasons: errs,
				Details: make([]interface{}, 0),
			},
		})
		return
	}

	var (
		feeds []*feedbacks.Feedback
	)
	feedDomain := feedbacks.NewFeedbackDomain(ds)
	filters := make(map[string]interface{})
	filters["distribution_id"] = payload.DistributionId
	feeds, _, err = feedDomain.FindAllWithFilter(req.GetContext().Value(), filters, true)

	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "SME:ERR:FND",
				Message: "Bad Request! Some problems occurred when searching the data.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	var factors *assessments.Template
	assessmentsDomain := assessments.NewAssessmentDomain(ds)
	factors, err = assessmentsDomain.FindTemplateDataByKey(req.GetContext().Value(), "review360")
	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "SMR:ERR:TPL",
				Message: "Internal error. Could not find factors.",
				Reasons: errs,
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	var resultItems []*FeedbackSummaryResponseItem
	resultItems, err = generateReviewSummaryData(feeds, factors, payload.ExcludeStatus...)
	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "SMR:ERR:GEN",
				Message: "Internal error. Could not transform data.",
				Reasons: errs,
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	writer := spreadsheets.NewExcelFile(fmt.Sprintf("exports/distribution_%d.%s", payload.DistributionId, payload.FileType), "Detail Summary", true)
	writer.Write(func(err error, results ...interface{}) {
		if err != nil {
			return
		}
		if results == nil || len(results) < 1 {
			return
		}
		var (
			wb *excelize.File
			ok bool
		)
		wb, ok = results[0].(*excelize.File)
		if !ok {
			return
		}
		writeReviewSummaryDetailSheet(writer, wb, resultItems)
		writeReviewSummaryBriefSheet(writer, wb, resultItems, "Summary")
		fmt.Println("done")
	})
	_ = vm.RenderJson(w, http.StatusOK, api.Response{
		Status: http.StatusOK,
		Content: map[string]interface{}{
			"status": "exporting...",
			"url":    "",
		},
	})
}

func writeReviewSummaryBriefSheet(writer spreadsheets.IWriter, wb *excelize.File, resultItems []*FeedbackSummaryResponseItem, sheet string) {
	sIdx := wb.GetSheetIndex(sheet)
	if sIdx < 0 {
		sIdx = wb.NewSheet(sheet)
	}
	wb.SetActiveSheet(sIdx)
	writer.SelectSheet(sheet)

	var (
		ax, col                      string
		rIdx, cIdx                   = 1, 1
		headerTableStyle, scoreStyle int
	)
	headerTableStyle, _ = wb.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 13.0},
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"FF9841"}},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	scoreStyle, _ = wb.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"FFDB71"}},
		Font: &excelize.Font{Bold: true, Size: 13},
	})
	headers := []string{"No", "Person", "Reviewers", "Status", "Score", "Rating"}
	colWidth := []float64{7, 22, 22, 10, 10, 18}
	for ih, header := range headers {
		ax = writer.GetAxis(rIdx, ih+1)
		col = writer.GetCol(ih + 1)
		_ = wb.SetCellStr(writer.GetSheet(), ax, header)
		_ = wb.SetCellStyle(writer.GetSheet(), ax, ax, headerTableStyle)
		_ = wb.SetColWidth(writer.GetSheet(), col, col, colWidth[ih])
	}
	for itemIdx, item := range resultItems {
		rIdx++
		cIdx = 1
		// Number
		_ = wb.SetCellInt(writer.GetSheet(), writer.GetAxis(rIdx, cIdx), itemIdx+1)
		// Person
		cIdx++
		_ = wb.SetCellStr(writer.GetSheet(), writer.GetAxis(rIdx, cIdx), item.Recipient.Name)
		if len(item.Items) < 1 {
			continue
		}
		// Feedback Status
		cIdx++
		for _, rItem := range item.Items {
			// Reviewer
			_ = wb.SetCellStr(writer.GetSheet(), writer.GetAxis(rIdx, cIdx), rItem.Respondent.Name)

			// Feedback Status
			_ = wb.SetCellStr(writer.GetSheet(), writer.GetAxis(rIdx, cIdx+1), strings.ToUpper(rItem.Status.String()))

			// Score
			_ = wb.SetCellValue(writer.GetSheet(), writer.GetAxis(rIdx, cIdx+2), rItem.Score)

			// Rating
			_ = wb.SetCellStr(writer.GetSheet(), writer.GetAxis(rIdx, cIdx+3), rItem.Rating)
			rIdx++
		}
		// Score Total
		_ = wb.SetCellValue(writer.GetSheet(), writer.GetAxis(rIdx, cIdx+2), fmt.Sprintf("%.02f", item.TotalScore))
		_ = wb.SetCellStr(writer.GetSheet(), writer.GetAxis(rIdx, cIdx+3), item.Rating)
		_ = wb.SetCellStyle(writer.GetSheet(), writer.GetAxis(rIdx, cIdx+2), writer.GetAxis(rIdx, cIdx+3), scoreStyle)
		rIdx += 2
	}
}

func writeReviewSummaryDetailSheet(writer spreadsheets.IWriter, wb *excelize.File, resultItems []*FeedbackSummaryResponseItem) {
	var (
		ax, col                                                    string
		rIdx, cIdx                                                 = 1, 1
		headerTableStyle, scoreStyle, strengthStyle, weaknessStyle int
	)
	headerTableStyle, _ = wb.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 13.0},
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"FF9841"}},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	scoreStyle, _ = wb.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"FFDB71"}},
		Font: &excelize.Font{Bold: true, Size: 13},
	})
	strengthStyle, _ = wb.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"7bCCf6"}},
	})
	weaknessStyle, _ = wb.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"f7F6ad"}},
	})
	headers := []string{"No", "Person", "Reviewers", "Status", "Topics", "SubTopic", "Weight", "Rating", "Notes"}
	colWidth := []float64{7, 22, 22, 10, 18, 42, 10, 10, 70}
	for ih, header := range headers {
		ax = writer.GetAxis(rIdx, ih+1)
		col = writer.GetCol(ih + 1)
		_ = wb.SetCellStr(writer.GetSheet(), ax, header)
		_ = wb.SetCellStyle(writer.GetSheet(), ax, ax, headerTableStyle)
		_ = wb.SetColWidth(writer.GetSheet(), col, col, colWidth[ih])
	}
	for itemIdx, item := range resultItems {
		rIdx++
		cIdx = 1
		// Number
		ax = writer.GetAxis(rIdx, cIdx)
		_ = wb.SetCellInt(writer.GetSheet(), ax, itemIdx+1)
		// Person
		cIdx++
		ax = writer.GetAxis(rIdx, cIdx)
		_ = wb.SetCellStr(writer.GetSheet(), ax, item.Recipient.Name)
		if len(item.Items) < 1 {
			continue
		}
		// Feedback Status
		cIdx++
		for _, rItem := range item.Items {
			// Reviewer
			ax = writer.GetAxis(rIdx, cIdx)
			_ = wb.SetCellStr(writer.GetSheet(), ax, rItem.Respondent.Name)
			// Feedback Status
			ax = writer.GetAxis(rIdx, cIdx+1)
			_ = wb.SetCellStr(writer.GetSheet(), ax, strings.ToUpper(rItem.Status.String()))
			// Topics
			// @todo: refactor this silly loop
			if rItem.Factors != nil {
				for _, factor := range rItem.Factors.Items {
					// Topic
					ax = writer.GetAxis(rIdx, cIdx+2)
					_ = wb.SetCellStr(writer.GetSheet(), ax, factor.Title)
					if len(factor.Items) > 0 {
						for _, factorItemLv2 := range factor.Items {
							if len(factorItemLv2.Items) < 1 {
								continue
							}
							for _, factorItemLv3 := range factorItemLv2.Items {
								if factorItemLv3.Rating < 1 {
									continue
								}
								// Sub Topic
								ax = writer.GetAxis(rIdx, cIdx+3)
								_ = wb.SetCellStr(writer.GetSheet(), ax, factorItemLv3.Title)
								// Weight
								ax = writer.GetAxis(rIdx, cIdx+4)
								_ = wb.SetCellValue(writer.GetSheet(), ax, fmt.Sprintf("%.0f%%", factorItemLv3.Weight))
								// Rating
								ax = writer.GetAxis(rIdx, cIdx+5)
								_ = wb.SetCellInt(writer.GetSheet(), ax, factorItemLv3.Rating)
								rIdx++
							}
						}
						continue
					}
					// Weight
					ax = writer.GetAxis(rIdx, cIdx+4)
					_ = wb.SetCellValue(writer.GetSheet(), ax, fmt.Sprintf("%.0f%%", factor.Weight))
					// Rating
					ax = writer.GetAxis(rIdx, cIdx+5)
					_ = wb.SetCellInt(writer.GetSheet(), ax, factor.Rating)
					rIdx++
				}
			}

			// Strength Topic
			ax = writer.GetAxis(rIdx, cIdx+2)
			_ = wb.SetCellStr(writer.GetSheet(), ax, "STRENGTH")
			// Strength Notes
			for _, strength := range rItem.Strengths {
				ax = writer.GetAxis(rIdx, cIdx+6)
				_ = wb.SetCellStyle(writer.GetSheet(), writer.GetAxis(rIdx, cIdx+2), writer.GetAxis(rIdx, cIdx+6), strengthStyle)
				_ = wb.SetCellStr(writer.GetSheet(), ax, strings.Title(strength))
				rIdx++
			}
			// Need Improvement Topic
			ax = writer.GetAxis(rIdx, cIdx+2)
			_ = wb.SetCellStr(writer.GetSheet(), ax, "NEED IMPROVEMENT")
			// Need Improvements Notes
			for _, needImprovement := range rItem.NeedImprovements {
				ax = writer.GetAxis(rIdx, cIdx+6)
				_ = wb.SetCellStyle(writer.GetSheet(), writer.GetAxis(rIdx, cIdx+2), writer.GetAxis(rIdx, cIdx+6), weaknessStyle)
				_ = wb.SetCellStr(writer.GetSheet(), ax, strings.Title(needImprovement))
				rIdx++
			}
			// Score
			rIdx++
			ax = writer.GetAxis(rIdx, cIdx+2)
			_ = wb.SetCellStr(writer.GetSheet(), ax, "SCORE")
			_ = wb.SetCellStyle(writer.GetSheet(), ax, writer.GetAxis(rIdx, cIdx+6), scoreStyle)
			_ = wb.SetCellValue(writer.GetSheet(), writer.GetAxis(rIdx, cIdx+5), rItem.Score)
			_ = wb.SetCellStr(writer.GetSheet(), writer.GetAxis(rIdx, cIdx+6), rItem.Rating)
			rIdx += 2
		}
	}
}
