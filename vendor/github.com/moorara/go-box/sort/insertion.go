package sort

import (
	. "github.com/moorara/go-box/dt"
)

// InsertionSort implements insertion sort algorithm
func InsertionSort(a []Generic, compare Compare) {
	n := len(a)
	for i := 0; i < n; i++ {
		for j := i; j > 0 && compare(a[j], a[j-1]) < 0; j-- {
			a[j], a[j-1] = a[j-1], a[j]
		}
	}
}
