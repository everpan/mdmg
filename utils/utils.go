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
	} else if i == len(v)-1 {
		return v[:i], ""
	}
	ver := v[i+1:]
	if len(ver) > 0 && ver[0] >= '0' && ver[0] <= '9' {
		return v[:i], ver
	}
	return v, ""
}
