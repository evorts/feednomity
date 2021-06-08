package utils

import "strconv"

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

func (a ArrayInt64) ToArrayString() []string {
	rs := make([]string, 0)
	for _, v := range a {
		rs = append(rs, strconv.FormatInt(v, 10))
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

type ArrayInterface []interface{}

func (a ArrayInterface) ToArrayString() ArrayString {
	rs := make(ArrayString, 0)
	for _, v := range a {
		if vv, ok := v.(string); ok {
			rs = append(rs, vv)
		}
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