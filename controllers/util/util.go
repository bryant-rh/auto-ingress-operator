package util

import (
	"strings"

	huge "github.com/dablelv/go-huge-util"
)

//IsValidServcieName 检查svc 名称前缀
func IsValidServcieName(name string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}

	return false
}

//IsValidNsName 检查namespace是否在配置内容中
func IsValidNsName(name string, namespaces interface{}) bool {
	m, _ := huge.ToMapSetE(namespaces)
	if _, ok := m[name]; ok {
		return true
	}
	return false
}
