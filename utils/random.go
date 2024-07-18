package utils

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	alphabets = "abcdefghijklmnopqrstuvwzyz"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Generate a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// Generate a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabets)

	for i := 0; i < n; i++ {
		c := alphabets[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomUser() string {
	return RandomString(6)
}

func RandomEmail() string {
	prefix := RandomString(6)
	return prefix + "@test.com"
}

func RandomAmount() string {
	n := RandomInt(100, 1000)
	amt := strconv.Itoa(int(n))
	return amt + ".00000000"
}

func RandomSmallAmount() string {
	n := RandomInt(1, 10)
	amt := strconv.Itoa(int(n))
	return amt + ".00000000"
}
