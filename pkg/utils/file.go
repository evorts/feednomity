package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func FileExist(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func ReadFile(filename string) string {
	return ReadAndReplace(filename, make(map[string]string))
}

func ReadAndReplace(filename string, placeholders map[string]string) string {
	if !FileExist(filename) {
		return ""
	}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return ""
	}
	rs := string(content)
	for k, v := range placeholders {
		rs = strings.ReplaceAll(rs, fmt.Sprintf("{{%s}}", k), v)
	}
	return rs
}