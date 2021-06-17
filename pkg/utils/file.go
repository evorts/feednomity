package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

type File string

func (f File) IsEmpty() bool {
	return len(f) < 1
}

func (f File) Value() string {
	return string(f)
}

func (f File) FilenameOnly() string {
	return filepath.Base(f.Value())
}

func (f File) FullPath() string {
	panic("implement me")
}

func (f File) IsExist() bool {
	_, err := os.Stat(f.Value())
	return !os.IsNotExist(err)
}

func (f File) Dir() string {
	currentDir, _ := os.Getwd()
	dir := filepath.Dir(f.Value())
	if len(dir) < 1 || !filepath.IsAbs(dir){
		dir = fmt.Sprintf("%s/%s", currentDir, dir)
	}
	return dir
}

func (f File) DirExist() bool {
	if _, err := os.Stat(f.Dir()); !os.IsNotExist(err) {
		return true
	}
	return false
}

func (f File) CreateDir() {
	if f.DirExist() {
		return
	}
	if err := os.Mkdir(f.Dir(), os.ModePerm); err != nil {
		fmt.Println(err)
	}
}

func (f File) InitFullPath() string {
	f.CreateDir()
	return fmt.Sprintf("%s/%s", f.Dir(), f.FilenameOnly())
}

type IFile interface {
	IsEmpty() bool
	Value() string
	FilenameOnly() string
	FullPath() string
	IsExist() bool
	Dir() string
	DirExist() bool
	CreateDir()
	InitFullPath() string
}
