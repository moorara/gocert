package ds

import (
	"testing"

	. "github.com/moorara/go-box/dt"
	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	tests := []struct {
		nodeSize             int
		compare              Compare
		enqueueItems         []string
		expectedSize         int
		expectedIsEmpty      bool
		expectedPeek         string
		expectedContains     []string
		expectedDequeueItems []string
	}{
		{
			2,
			CompareString,
			[]string{},
			0, true,
			"",
			[]string{},
			[]string{},
		},
		{
			2,
			CompareString,
			[]string{"a", "b"},
			2, false,
			"a",
			[]string{"a", "b"},
			[]string{"a", "b"},
		},
		{
			2,
			CompareString,
			[]string{"a", "b", "c"},
			3, false,
			"a",
			[]string{"a", "b", "c"},
			[]string{"a", "b", "c"},
		},
		{
			2,
			CompareString,
			[]string{"a", "b", "c", "d", "e", "f", "g"},
			7, false,
			"a",
			[]string{"a", "b", "c", "d", "e", "f", "g"},
			[]string{"a", "b", "c", "d", "e", "f", "g"},
		},
	}

	for _, test := range tests {
		queue := NewQueue(test.nodeSize, test.compare)

		// Queue initially should be empty
		assert.Zero(t, queue.Size())
		assert.True(t, queue.IsEmpty())
		assert.Nil(t, queue.Peek())
		queue.Contains(nil)
		assert.Nil(t, queue.Dequeue())

		for _, item := range test.enqueueItems {
			queue.Enqueue(item)
		}

		assert.Equal(t, test.expectedSize, queue.Size())
		assert.Equal(t, test.expectedIsEmpty, queue.IsEmpty())

		if test.expectedSize == 0 {
			assert.Nil(t, queue.Peek())
		} else {
			assert.Equal(t, test.expectedPeek, queue.Peek())
		}

		for _, item := range test.expectedContains {
			assert.True(t, queue.Contains(item))
		}

		for _, item := range test.expectedDequeueItems {
			assert.Equal(t, item, queue.Dequeue())
		}

		// Queue should be empty at the end
		assert.Zero(t, queue.Size())
		assert.True(t, queue.IsEmpty())
		assert.Nil(t, queue.Peek())
		queue.Contains(nil)
		assert.Nil(t, queue.Dequeue())
	}
}

func BenchmarkQueue(b *testing.B) {
	nodeSize := 1024
	item := 27

	b.Run("Enqueue", func(b *testing.B) {
		queue := NewQueue(nodeSize, CompareInt)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			queue.Enqueue(item)
		}
	})

	b.Run("Dequeue", func(b *testing.B) {
		queue := NewQueue(nodeSize, CompareInt)
		for n := 0; n < b.N; n++ {
			queue.Enqueue(item)
		}
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			queue.Dequeue()
		}
	})
}
