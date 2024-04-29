package utils 

import "strings"

// Join concatenates mulitiple strings
func StrJoin(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}
