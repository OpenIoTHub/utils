package stringUtil

import (
	"regexp"
	"strconv"
)

func ConvertOctonaryUtf8(in string) string {
	s := []byte(in)
	reg := regexp.MustCompile(`\\[0-9]{3}`)

	out := reg.ReplaceAllFunc(s,
		func(b []byte) []byte {
			i, _ := strconv.ParseInt(string(b[1:]), 10, 0)
			return []byte{byte(i)}
		})
	return string(out)
}
