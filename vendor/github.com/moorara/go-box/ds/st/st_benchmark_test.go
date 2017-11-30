package st

import (
	"math/rand"
	"testing"

	. "github.com/moorara/go-box/dt"
	"github.com/moorara/go-box/util"
)

const (
	seed = 27
)

func getIntSlice(size int) []Generic {
	items := make([]Generic, size)
	for i := 0; i < len(items); i++ {
		items[i] = i
	}
	util.Shuffle(items)

	return items
}

func runPutBenchmark(b *testing.B, ost OrderedSymbolTable) {
	items := getIntSlice(b.N)
	rand.Seed(seed)
	util.Shuffle(items)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		ost.Put(items[n], "")
	}
}

func runGetBenchmark(b *testing.B, ost OrderedSymbolTable) {
	items := getIntSlice(b.N)
	rand.Seed(seed)
	util.Shuffle(items)
	for n := 0; n < b.N; n++ {
		ost.Put(items[n], "")
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		ost.Get(items[n])
	}
}

func BenchmarkOrderedSymbolTable(b *testing.B) {
	b.Run("BST.Put", func(b *testing.B) {
		ost := NewBST(CompareInt)
		runPutBenchmark(b, ost)
	})

	b.Run("BST.Get", func(b *testing.B) {
		ost := NewBST(CompareInt)
		runGetBenchmark(b, ost)
	})

	b.Run("AVL.Put", func(b *testing.B) {
		ost := NewAVL(CompareInt)
		runPutBenchmark(b, ost)
	})

	b.Run("AVL.Get", func(b *testing.B) {
		ost := NewAVL(CompareInt)
		runGetBenchmark(b, ost)
	})

	b.Run("RedBlack.Put", func(b *testing.B) {
		ost := NewRedBlack(CompareInt)
		runPutBenchmark(b, ost)
	})

	b.Run("RedBlack.Get", func(b *testing.B) {
		ost := NewRedBlack(CompareInt)
		runGetBenchmark(b, ost)
	})
}
