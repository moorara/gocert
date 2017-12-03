package util

import (
	. "github.com/moorara/go-box/dt"
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
