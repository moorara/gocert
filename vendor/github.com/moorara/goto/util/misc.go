package util

import (
	. "github.com/moorara/goto/dt"
)

// IsSorted checks if an array is ascendingly sorted
func IsSorted(items []Generic, compare Compare) bool {
	for i := 0; i < len(items)-1; i++ {
		if compare(items[i], items[i+1]) > 0 {
			return false
		}
	}

	return true
}

// IsStringIn checks if a string is in a list of strings
func IsStringIn(s string, list ...string) bool {
	for _, str := range list {
		if str == s {
			return true
		}
	}

	return false
}
