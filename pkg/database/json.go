package database

import "database/sql/driver"

type IJson interface {
	Value () (value driver.Value, err error)
	Scan (value interface{}) error
}
