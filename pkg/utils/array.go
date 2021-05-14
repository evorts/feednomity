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

type ArrayString []string

func (a ArrayString) Reduce() []string {
	rs := make([]string, 0)
	for _, v := range a {
		if len(v) < 1 {
			continue
		}
		rs = append(rs, v)
	}
	return rs
}

func (a ArrayString) ToArrayInterface() []interface{} {
	rs := make([]interface{}, 0)
	for _, v := range a {
		rs = append(rs, v)
	}
	return rs
}


func InArray(arr []interface{}, v interface{}) bool  {
	for _, av := range arr {
		if av == v {
			return true
		}
	}
	return false
}