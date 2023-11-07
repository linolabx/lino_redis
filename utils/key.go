package utils

import (
	"strings"
)

func KeyJoin(pathList ...string) string {
	parts := []string{}

	for _, p := range pathList {
		if p != "" {
			parts = append(parts, strings.Trim(p, ":"))
		}
	}
	return strings.Join(parts, ":")
}

func KeyJoinSlice(path string, pathList ...string) string {
	return KeyJoin(append([]string{path}, pathList...)...)
}
