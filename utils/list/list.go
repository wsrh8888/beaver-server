package utils

import (
	"fmt"
	"regexp"

	"github.com/zeromicro/go-zero/core/logx"
)

func InListByRegex(list []string, key string) (ok bool) {

	for _, s := range list {
		regex, err := regexp.Compile(s)
		if err != nil {
			logx.Errorf("compile regex error: %v", err)
			return
		}
		fmt.Println(key, s)
		if regex.MatchString(key) {
			return true
		}
	}
	return false
}

func InList(list []string, key string) bool {
	for _, i := range list {
		if i == key {
			return true
		}
	}
	return false
}
