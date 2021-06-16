package utils

import "strings"

func ParseTag(tagField string) map[string]string {
	res := make(map[string]string, 0)
	if len(tagField) < 1 || !strings.Contains(tagField, ":") {
		return res
	}
	tags := strings.Split(tagField, " ")
	for _, t := range tags {
		kv := strings.Split(t, ":")
		if len(kv) != 2 {
			continue
		}
		res[kv[0]] = kv[1]
	}
	return res
}

func GetTagName(key, tagField string) string {
	tags := ParseTag(tagField)
	if v, ok := tags[key]; ok {
		return strings.Trim(v, "\"")
	}
	return ""
}
