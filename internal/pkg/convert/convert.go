package convert

import (
	"strconv"
	"strings"
)

func Uint32SliceToString(slice []uint32, sep string) string {
	var sb strings.Builder
	for i, num := range slice {
		if i > 0 {
			sb.WriteString(sep)
		}
		sb.WriteString(strconv.Itoa(int(num)))
	}
	if len(slice) == 0 {
		sb.WriteString("-")
	}
	return sb.String()
}
