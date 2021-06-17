package spreadsheets

type IWriter interface {
	Write(callback func(err error, results ...interface{}))
	GetAxis(row, col int) string
	GetCol(col int) string
	GetSheet() string
	SelectSheet(name string)
}
