package sort

import (
	. "github.com/moorara/go-box/dt"
)

func sink(a []Generic, k, n int, compare Compare) {
	for 2*k <= n {
		j := 2 * k
		if j < n && compare(a[j], a[j+1]) < 0 {
			j++
		}
		if compare(a[k], a[j]) >= 0 {
			break
		}
		a[k], a[j] = a[j], a[k]
		k = j
	}
}

func heapSort(a []Generic, compare Compare) {
	n := len(a) - 1

	for k := n / 2; k >= 1; k-- { // build max-heap bottom-up
		sink(a, k, n, compare)
	}
	for n > 1 { // remove the maximum, one at a time
		a[1], a[n] = a[n], a[1]
		n--
		sink(a, 1, n, compare)
	}
}

// HeapSort implements heap sort algorithm
func HeapSort(a []Generic, compare Compare) {
	// Heap elements need to start from position 1
	aux := append([]Generic{nil}, a...)
	heapSort(aux, compare)
	copy(a, aux[1:])
}
