package utils

type Page int

func (p Page) ToOffset(limit int) int {
	return(p.ToInt() - 1) * limit
}

func (p Page) ToInt() int {
	if p < 1 {
		return 1
	}
	return int(p)
}
