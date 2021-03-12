package utils

type ArrayInteger []int

func (a ArrayInteger) ToArrayInterface() []interface{} {
	rs := make([]interface{}, 0)
	for _, v := range a {
		rs = append(rs, v)
	}
	return rs
}

type ArrayInt64 []int64

func (a ArrayInt64) ToArrayInterface() []interface{} {
	rs := make([]interface{}, 0)
	for _, v := range a {
		rs = append(rs, v)
	}
	return rs
}