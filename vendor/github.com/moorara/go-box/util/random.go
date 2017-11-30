package util

import (
	"math/rand"
	"time"

	. "github.com/moorara/go-box/dt"
)

// SeedWithNow seeds the random generator with time now
func SeedWithNow() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Shuffle shuffles an array in O(n) time
func Shuffle(a []Generic) {
	n := len(a)
	for i := 0; i < n; i++ {
		r := i + rand.Intn(n-i)
		a[i], a[r] = a[r], a[i]
	}
}

// GenerateInt generates a random integer
func GenerateInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

// GenerateString generates a random string
func GenerateString(minLen, maxLen int) string {
	strLen := GenerateInt(minLen, maxLen)
	bytes := make([]byte, strLen)
	for j := 0; j < strLen; j++ {
		bytes[j] = byte(GenerateInt(65, 95))
	}

	return string(bytes)
}

// GenerateIntSlice generates an array with random integers
func GenerateIntSlice(size, min, max int) []Generic {
	items := make([]Generic, size)
	for i := 0; i < len(items); i++ {
		items[i] = GenerateInt(min, max)
	}

	return items
}

// GenerateStringSlice generates an array with random strings
func GenerateStringSlice(size, minLen, maxLen int) []Generic {
	items := make([]Generic, size)
	for i := 0; i < len(items); i++ {
		items[i] = GenerateString(minLen, maxLen)
	}

	return items
}
