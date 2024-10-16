package gee

import "strings"

var builder strings.Builder

func joinStr(strs ...string) string {
	builder.Reset()
	for i := range strs {
		builder.WriteString(strs[i])
	}
	return builder.String()
}
