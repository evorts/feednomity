package mailer

import (
	"fmt"
	"github.com/evorts/feednomity/pkg/utils"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
)

func readTemplate(path string) string {
	if !utils.FileExist(path) {
		return ""
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(content)
}

func bindDataToTemplate(data map[string]string, template string) string {
	if len(data) < 1 {
		return template
	}
	for k, v := range data {
		template = strings.ReplaceAll(template, fmt.Sprintf("{{%s}}",k), v)
	}
	return template
}

func validate(to []Target, subject, content string) error {
	if len(to) < 1 {
		return errors.New("missing recipients data!")
	}
	if len(subject) < 1 {
		return errors.New("missing subject!")
	}
	if len(content) < 1 {
		return errors.New("missing content data")
	}
	return nil
}