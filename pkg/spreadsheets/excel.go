package spreadsheets

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/evorts/feednomity/pkg/utils"
	"github.com/pkg/errors"
	"math"
	"os"
)

var alphabet = []string{
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
	"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
}

type excel struct {
	filename       string
	sheet          string
	replaceIfExist bool
	wb             *excelize.File
}

func (e *excel) GetAxis(row, col int) string {
	col -= 1
	if col < 0 {
		col = 0
	}
	colMultiplier := int(math.Floor(float64(col) / float64(len(alphabet))))
	colMod := col % len(alphabet)
	if colMultiplier < 1 {
		return fmt.Sprintf("%s%d", alphabet[colMod], row)
	}
	return fmt.Sprintf("%s%s%d", alphabet[colMultiplier-1], alphabet[colMod], row)
}

func (e *excel) GetCol(col int) string {
	col -= 1
	if col < 0 {
		col = 0
	}
	colMultiplier := int(math.Floor(float64(col) / float64(len(alphabet))))
	colMod := col % len(alphabet)
	if colMultiplier < 1 {
		return fmt.Sprintf("%s", alphabet[colMod])
	}
	return fmt.Sprintf("%s%s", alphabet[colMultiplier-1], alphabet[colMod])
}

func (e *excel) Write(callback func(err error, results ...interface{})) {
	var err error
	file := utils.File(e.filename)
	if !e.replaceIfExist && file.IsExist() {
		callback(errors.New("file already exist. process stopped."))
		return
	}
	if e.replaceIfExist && file.IsExist() {
		if err = os.Remove(e.filename); err != nil {
			callback(err)
			return
		}
	}
	f := file.InitFullPath()
	e.wb, err = excelize.OpenFile(f)
	if err != nil {
		fmt.Println(fmt.Errorf("unable to open excel file. error: %w", err))
		e.wb = excelize.NewFile()
		e.wb.Path = f
	}
	sIdx := e.wb.GetSheetIndex(e.GetSheet())
	if sIdx < 0 {
		sIdx = e.wb.NewSheet(e.GetSheet())
	}
	e.wb.SetActiveSheet(sIdx)
	callback(nil, e.wb)
	err = e.wb.Save()
	fmt.Println(err)
}

func (e *excel) GetSheet() string {
	return e.sheet
}

func (e *excel) SelectSheet(name string) {
	sIdx := e.wb.GetSheetIndex(name)
	if sIdx < 0 {
		sIdx = e.wb.NewSheet(name)
	}
	e.wb.SetActiveSheet(sIdx)
	e.sheet = name
}

func NewExcelFile(filename string, sheet string, replace bool) IWriter {
	return &excel{filename: filename, sheet: sheet, replaceIfExist: replace}
}
