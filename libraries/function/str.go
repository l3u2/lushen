package function

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"
)

func Implode(glue string, pieces []string) string {
	var buf bytes.Buffer
	l := len(pieces)
	for _, str := range pieces {
		buf.WriteString(str)
		if l--; l > 0 {
			buf.WriteString(glue)
		}
	}
	return buf.String()
}

func Md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

func Explode(delimiter, str string) []string {
	return strings.Split(str, delimiter)
}

func StrToInt(str string) int {
	ret, _ := strconv.Atoi(str)
	return ret
}

func StrReplace(search, replace, subject string, count int) string {
	return strings.Replace(subject, search, replace, count)
}
