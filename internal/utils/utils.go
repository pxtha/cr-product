package utils

import (
	"log"
	"math/rand"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func RandomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(RandIntRange(65, 90))
	}
	return string(bytes)
}

func RandIntRange(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func FMPrice(price string) string {
	rs := strings.Replace(price, " ₫", "", -1)
	rs = strings.Replace(rs, ".", "", -1)
	return rs
}

func CheckError(err error) {
	if err != nil {
		log.Println(err)
	}
}
func NormalizeString(str string) string {
	trans := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(trans, str)
	result = strings.ReplaceAll(result, "đ", "d")
	result = strings.ReplaceAll(result, "Đ", "D")
	return result
}
