package handler

type Page int

func (p Page) Value() int {
	if p < 1 {
		return 1
	}
	return int(p)
}

type Limit int

func (l Limit) Value() int {
	if l < 1 {
		return 10
	}
	return int(l)
}

