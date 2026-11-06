package ymsgw

import "strings"

// trim 去除首尾空白。
func trim(s string) string { return strings.TrimSpace(s) }

// splitLines 按 \n 切段（去掉 \r），保留段落顺序。
func splitLines(s string) []string {
	s = strings.ReplaceAll(s, "\r", "")
	return strings.Split(s, "\n")
}

// roleRank 角色序，与 doc_rpc 权限模型对齐：owner>manager>editor>commenter>viewer。
func roleRank(role string) int {
	switch role {
	case "owner":
		return 5
	case "manager":
		return 4
	case "editor":
		return 3
	case "commenter":
		return 2
	case "viewer":
		return 1
	default:
		return 0
	}
}
