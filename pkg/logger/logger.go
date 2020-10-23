package logger

import (
	"encoding/json"
	"fmt"
	"log"
)

type logField struct {
	Component string `json:"component"`
	Content interface{} `json:"content"`
}

type logger struct{}

type IManager interface {
	Log(key string, value interface{})
	Fatal(value interface{})
}

func NewLogger() IManager {
	return &logger{}
}

func (l logger) Log(key string, value interface{}) {
	if j, err := json.Marshal(logField{ key, value}); err == nil {
		fmt.Println(string(j))
	}
}

func (l logger) Fatal(value interface{}) {
	j, _ := json.Marshal(logField{ "fatal", value})
	log.Fatal(j)
}