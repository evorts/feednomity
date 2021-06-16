package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func GenerateFilters(values map[string]interface{}) (sql string, args []interface{}) {
	arr := make([]string, 0)
	for k, v := range values {
		rt := reflect.TypeOf(v)
		switch rt.Kind() {
		case reflect.Slice, reflect.Array:
			vArrS, okS := v.([]string)
			if okS {
				arr = append(
					arr,
					fmt.Sprintf(
						"%s IN (%s)", k,
						strings.TrimRight(strings.Repeat("?,", len(vArrS)), ", "),
					),
				)
				args = append(args, ArrayString(vArrS).ToArrayInterface()...)
				continue
			}
			vArrI64, okI64 := v.([]int64)
			if okI64 {
				arr = append(
					arr,
					fmt.Sprintf(
						"%s IN (%s)", k,
						strings.TrimRight(strings.Repeat("?,", len(vArrI64)), ", "),
					),
				)
				args = append(args, ArrayInt64(vArrI64).ToArrayInterface()...)
				continue
			}
			vArrInt, okInt := v.([]int)
			if okInt {
				arr = append(
					arr,
					fmt.Sprintf(
						"%s IN (%s)", k,
						strings.TrimRight(strings.Repeat("?,", len(vArrInt)), ", "),
					),
				)
				args = append(args, ArrayInteger(vArrInt).ToArrayInterface()...)
				continue
			}
		default:
			arr = append(arr, fmt.Sprintf("%s = ?", k))
			args = append(args, v)
		}
	}
	sql = strings.Join(arr, " AND ")
	return
}
