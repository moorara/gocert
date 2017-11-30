package dt

import (
	"strconv"
	"strings"
)

// CompareInt compares two int numbers
func CompareInt(a Generic, b Generic) int {
	intA, _ := a.(int)
	intB, _ := b.(int)
	diff := intA - intB
	switch {
	case diff < 0:
		return -1
	case diff > 0:
		return 1
	default:
		return 0
	}
}

// BitStringInt returns the bit-string representation of an int number
func BitStringInt(a Generic) []byte {
	intA, _ := a.(int)
	return []byte(strconv.Itoa(intA))
}

// CompareString compares two strings
func CompareString(a Generic, b Generic) int {
	strA, _ := a.(string)
	strB, _ := b.(string)
	return strings.Compare(strA, strB)
}

// BitStringString returns the bit-string representation of a string
func BitStringString(a Generic) []byte {
	strA, _ := a.(string)
	return []byte(strA)
}
