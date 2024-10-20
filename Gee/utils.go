package gee

import "strings"

const (
	urlSep = "/"
)

// 将多个字符串进行拼接
func joinStr(strs ...string) string {
	var builder strings.Builder
	for i := range strs {
		builder.WriteString(strs[i])
	}
	return builder.String()
}

// 解析pattern
// - 遇到*前缀的通配符直接停止
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, urlSep)

	parts := make([]string, 0)
	for i := range vs {
		if vs[i] != "" {
			parts = append(parts, vs[i])
			if vs[i][0] == '*' {
				break
			}
		}
	}
	return parts
}

// 按照sep分隔符，分割str并返回切分得到的
func splitStr(str string, sep string) []string {
	strs := strings.Split(str, sep)
	var result []string
	for i := range strs { // 过滤掉空字符串
		if strs[i] != "" {
			result = append(result, strs[i])
		}
	}
	return result
}
