package heap

import (
	"testing"

	. "github.com/moorara/go-box/dt"
	"github.com/moorara/go-box/util"
	"github.com/stretchr/testify/assert"
)

func TestMinHeap(t *testing.T) {
	tests := []struct {
		initialSize           int
		compareKey            Compare
		compareValue          Compare
		insertKeys            []int
		insertValues          []string
		expectedSize          int
		expectedIsEmpty       bool
		expectedPeekKey       int
		expectedPeekValue     string
		expectedContainsKey   []int
		expectedContainsValue []string
		expectedDeleteKeys    []int
		expectedDeleteValues  []string
	}{
		{
			2,
			CompareInt, CompareString,
			[]int{}, []string{},
			0, true,
			0, "",
			[]int{}, []string{},
			[]int{}, []string{},
		},
		{
			2,
			CompareInt, CompareString,
			[]int{30, 10, 20}, []string{"thirty", "ten", "twenty"},
			3, false,
			10, "ten",
			[]int{10, 20, 30}, []string{"ten", "twenty", "thirty"},
			[]int{10, 20, 30}, []string{"ten", "twenty", "thirty"},
		},
		{
			4,
			CompareInt, CompareString,
			[]int{50, 30, 40, 10, 20}, []string{"fifty", "thirty", "forty", "ten", "twenty"},
			5, false,
			10, "ten",
			[]int{10, 20, 30, 40, 50}, []string{"ten", "twenty", "thirty", "forty", "fifty"},
			[]int{10, 20, 30, 40, 50}, []string{"ten", "twenty", "thirty", "forty", "fifty"},
		},
		{
			4,
			CompareInt, CompareString,
			[]int{90, 80, 70, 40, 50, 60, 30, 10, 20}, []string{"ninety", "eighty", "seventy", "forty", "fifty", "sixty", "thirty", "ten", "twenty"},
			9, false,
			10, "ten",
			[]int{10, 20, 30, 40, 50, 60, 70, 80, 90}, []string{"ten", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety"},
			[]int{10, 20, 30, 40, 50, 60, 70, 80, 90}, []string{"ten", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety"},
		},
	}

	for _, test := range tests {
		heap := NewMinHeap(test.initialSize, test.compareKey, test.compareValue)

		// Heap initially should be empty
		peekKey, peekValue := heap.Peek()
		deleteKey, deleteValue := heap.Delete()
		assert.Nil(t, peekKey)
		assert.Nil(t, peekValue)
		assert.Nil(t, deleteKey)
		assert.Nil(t, deleteValue)
		assert.Zero(t, heap.Size())
		assert.True(t, heap.IsEmpty())
		assert.False(t, heap.ContainsKey(nil))
		assert.False(t, heap.ContainsValue(nil))

		for i := 0; i < len(test.insertKeys); i++ {
			heap.Insert(test.insertKeys[i], test.insertValues[i])
		}

		assert.Equal(t, test.expectedSize, heap.Size())
		assert.Equal(t, test.expectedIsEmpty, heap.IsEmpty())

		peekKey, peekValue = heap.Peek()
		if test.expectedSize == 0 {
			assert.Nil(t, peekKey)
			assert.Nil(t, peekValue)
		} else {
			assert.Equal(t, test.expectedPeekKey, peekKey)
			assert.Equal(t, test.expectedPeekValue, peekValue)
		}

		for _, key := range test.expectedContainsKey {
			assert.True(t, heap.ContainsKey(key))
		}

		for _, value := range test.expectedContainsValue {
			assert.True(t, heap.ContainsValue(value))
		}

		for i := 0; i < len(test.expectedDeleteKeys); i++ {
			deleteKey, deleteValue = heap.Delete()
			assert.Equal(t, test.expectedDeleteKeys[i], deleteKey)
			assert.Equal(t, test.expectedDeleteValues[i], deleteValue)
		}

		// Heap should be empty at the end
		peekKey, peekValue = heap.Peek()
		deleteKey, deleteValue = heap.Delete()
		assert.Nil(t, peekKey)
		assert.Nil(t, peekValue)
		assert.Nil(t, deleteKey)
		assert.Nil(t, deleteValue)
		assert.Zero(t, heap.Size())
		assert.True(t, heap.IsEmpty())
		assert.False(t, heap.ContainsKey(nil))
		assert.False(t, heap.ContainsValue(nil))
	}
}

func BenchmarkMinHeap(b *testing.B) {
	heapSize := 1024
	minInt := 0
	maxInt := 1000000
	util.SeedWithNow()

	b.Run("Insert", func(b *testing.B) {
		heap := NewMinHeap(heapSize, CompareInt, CompareString)
		items := util.GenerateIntSlice(b.N, minInt, maxInt)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			heap.Insert(items[n], "")
		}
	})

	b.Run("Delete", func(b *testing.B) {
		heap := NewMinHeap(heapSize, CompareInt, CompareString)
		items := util.GenerateIntSlice(b.N, minInt, maxInt)
		for n := 0; n < b.N; n++ {
			heap.Insert(items[n], "")
		}
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			heap.Delete()
		}
	})
}
