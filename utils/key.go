package utils

import (
	"strings"
)

func KeyJoin(pathList ...string) string {
	parts := []string{}

	for _, p := range pathList {
		parts = append(parts, strings.Trim(p, ":"))
	}
	return strings.Join(parts, ":")
}

func KeyJoinSlice(path string, pathList ...string) string {
	if len(pathList) == 0 {
		return strings.Trim(path, ":")
	}
	return strings.Trim(path, ":") + ":" + KeyJoin(pathList...)
}
