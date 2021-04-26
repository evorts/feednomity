package utils

import (
	"errors"
	"reflect"
)

func MergeStruct(dst interface{}, src interface{}, skipFieldsOnEmpty []string) error {
	var v1, v2 reflect.Value
	var t1, t2 reflect.Type

	v1 = reflect.ValueOf(dst)
	v2 = reflect.ValueOf(src)
	t1 = reflect.TypeOf(dst)
	t2 = reflect.TypeOf(src)

	if v1.Kind() == reflect.Ptr {
		t1 = v1.Elem().Type()
		v1 = reflect.Indirect(v1)
	}
	if v2.Kind() == reflect.Ptr {
		t2 = v2.Elem().Type()
		v2 = reflect.Indirect(v2)
	}
	if t1.Kind() != reflect.Struct || t2.Kind() != reflect.Struct || !reflect.DeepEqual(t1, t2) {
		return errors.New("invalid arguments data type")
	}
	for i := 0; i < v1.NumField(); i++ {
		fn1 := t1.Field(i).Name
		srcLoop:
		for j := 0; j < v2.NumField(); j++ {
			if fn1 != t2.Field(j).Name {
				continue
			}
			fv := v2.Field(j)
			if len(skipFieldsOnEmpty) > 0 {
				for _, v := range skipFieldsOnEmpty {
					if v != fn1 {
						continue
					}
					switch fv.Kind() {
					case reflect.Ptr, reflect.Interface:
						if fv.IsNil() {
							break srcLoop
						}
					case reflect.String:
						if fv.String() == "" {
							break srcLoop
						}
					default:
					}
					if fv.Kind() == reflect.Ptr && fv.Kind() == reflect.Interface {
						if fv.IsNil() {
							break srcLoop
						}
					} else if fv.IsZero() {
						break srcLoop
					}
				}
			}
			switch fv.Kind() {
			case reflect.Ptr, reflect.Interface:
				if fv.IsNil() {
					break
				}
			default:
				v1.Field(i).Set(v2.Field(j))
			}
		}
	}
	return nil
}

func TransformStruct(dst interface{}, src interface{}) error {
	var v1, v2 reflect.Value
	var t1, t2 reflect.Type

	v1 = reflect.ValueOf(dst)
	v2 = reflect.ValueOf(src)
	t1 = reflect.TypeOf(dst)
	t2 = reflect.TypeOf(src)

	if v1.Kind() == reflect.Ptr {
		t1 = v1.Elem().Type()
		v1 = reflect.Indirect(v1)
	}
	if v2.Kind() == reflect.Ptr {
		t2 = v2.Elem().Type()
		v2 = reflect.Indirect(v2)
	}
	if t1.Kind() != reflect.Struct || t2.Kind() != reflect.Struct {
		return errors.New("invalid arguments data type")
	}
	for i := 0; i < v1.NumField(); i++ {
		fn1 := t1.Field(i).Name
		for j := 0; j < v2.NumField(); j++ {
			if fn1 != t2.Field(j).Name {
				continue
			}
			fv := v2.Field(j)
			switch fv.Kind() {
			case reflect.Ptr, reflect.Interface:
				if fv.IsNil() {
					break
				}
			default:
				v1.Field(i).Set(v2.Field(j))
			}
		}
	}
	return nil
}

