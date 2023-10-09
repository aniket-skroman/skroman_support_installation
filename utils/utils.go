package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func RandomString(lenght int) (ans string) {
	chars := "abcdefghijklmnopqrstuvwxyz"
	for len(ans) < lenght {
		ans += string(chars[rand.Intn(len(chars))])
	}

	return
}

func RandomInt(length int) (ans string) {
	t := fmt.Sprint(time.Now().Nanosecond())
	ans = t[:7]
	return
}
