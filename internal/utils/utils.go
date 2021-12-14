package utils

import (
	"log"
	"math/rand"
	"strings"
	"time"
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
