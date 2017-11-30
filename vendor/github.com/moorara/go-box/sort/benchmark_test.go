package sort

import (
	"math/rand"
	"sort"
	"testing"

	. "github.com/moorara/go-box/dt"
	"github.com/moorara/go-box/util"
)

const (
	seed   = 27
	size   = 1000
	minInt = 0
	maxInt = 1000000
)

type GenricSlice []Generic

func (s GenricSlice) Len() int {
	return len(s)
}

func (s GenricSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s GenricSlice) Less(i, j int) bool {
	return CompareInt(s[i], s[j]) < 0
}

func BenchmarkSort(b *testing.B) {
	b.Run("sort.Sort", func(b *testing.B) {
		rand.Seed(seed)
		items := GenricSlice(util.GenerateIntSlice(size, minInt, maxInt))
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			util.Shuffle(items)
			sort.Sort(items)
		}
	})

	b.Run("HeapSort", func(b *testing.B) {
		rand.Seed(seed)
		items := util.GenerateIntSlice(size, minInt, maxInt)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			util.Shuffle(items)
			HeapSort(items, CompareInt)
		}
	})

	b.Run("InsertionSort", func(b *testing.B) {
		rand.Seed(seed)
		items := util.GenerateIntSlice(size, minInt, maxInt)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			util.Shuffle(items)
			InsertionSort(items, CompareInt)
		}
	})

	b.Run("MergeSort", func(b *testing.B) {
		rand.Seed(seed)
		items := util.GenerateIntSlice(size, minInt, maxInt)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			util.Shuffle(items)
			MergeSort(items, CompareInt)
		}
	})

	b.Run("MergeSortRec", func(b *testing.B) {
		rand.Seed(seed)
		items := util.GenerateIntSlice(size, minInt, maxInt)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			util.Shuffle(items)
			MergeSortRec(items, CompareInt)
		}
	})

	b.Run("QuickSort", func(b *testing.B) {
		rand.Seed(seed)
		items := util.GenerateIntSlice(size, minInt, maxInt)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			util.Shuffle(items)
			QuickSort(items, CompareInt)
		}
	})

	b.Run("QuickSort3Way", func(b *testing.B) {
		rand.Seed(seed)
		items := util.GenerateIntSlice(size, minInt, maxInt)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			util.Shuffle(items)
			QuickSort3Way(items, CompareInt)
		}
	})

	b.Run("ShellSort", func(b *testing.B) {
		rand.Seed(seed)
		items := util.GenerateIntSlice(size, minInt, maxInt)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			util.Shuffle(items)
			ShellSort(items, CompareInt)
		}
	})
}
