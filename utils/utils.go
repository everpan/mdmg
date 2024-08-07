package utils

import "strings"

func SplitModuleVersion(v string) (string, string) {
	// v := strings.TrimSpace(val)
	if len(v) == 0 {
		return "", ""
	}
	i := strings.LastIndex(v, "-")
	if i == -1 {
		return v, ""
	}
	return v[:i], v[i+1:]
}
