package ds

import (
	"testing"

	. "github.com/moorara/go-box/dt"
	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	tests := []struct {
		nodeSize         int
		compare          Compare
		pushItems        []string
		expectedSize     int
		expectedIsEmpty  bool
		expectedPeek     string
		expectedContains []string
		expectedPopItems []string
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
			"b",
			[]string{"a", "b"},
			[]string{"b", "a"},
		},
		{
			2,
			CompareString,
			[]string{"a", "b", "c"},
			3, false,
			"c",
			[]string{"a", "b", "c"},
			[]string{"c", "b", "a"},
		},
		{
			2,
			CompareString,
			[]string{"a", "b", "c", "d", "e", "f", "g"},
			7, false,
			"g",
			[]string{"a", "b", "c", "d", "e", "f", "g"},
			[]string{"g", "f", "e", "d", "c", "b", "a"},
		},
	}

	for _, test := range tests {
		stack := NewStack(test.nodeSize, test.compare)

		// Stack initially should be empty
		assert.Zero(t, stack.Size())
		assert.True(t, stack.IsEmpty())
		assert.Nil(t, stack.Pop())
		assert.Nil(t, stack.Peek())
		assert.False(t, stack.Contains(nil))

		for _, item := range test.pushItems {
			stack.Push(item)
		}

		assert.Equal(t, test.expectedSize, stack.Size())
		assert.Equal(t, test.expectedIsEmpty, stack.IsEmpty())

		if test.expectedSize == 0 {
			assert.Nil(t, stack.Peek())
		} else {
			assert.Equal(t, test.expectedPeek, stack.Peek())
		}

		for _, item := range test.expectedContains {
			assert.True(t, stack.Contains(item))
		}

		for _, item := range test.expectedPopItems {
			assert.Equal(t, item, stack.Pop())
		}

		// Stack should be empty at the end
		assert.Zero(t, stack.Size())
		assert.True(t, stack.IsEmpty())
		assert.Nil(t, stack.Pop())
		assert.Nil(t, stack.Peek())
		assert.False(t, stack.Contains(nil))
	}
}

func BenchmarkStack(b *testing.B) {
	nodeSize := 1024
	item := 27

	b.Run("Push", func(b *testing.B) {
		stack := NewStack(nodeSize, CompareInt)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			stack.Push(item)
		}
	})

	b.Run("Pop", func(b *testing.B) {
		stack := NewStack(nodeSize, CompareInt)
		for n := 0; n < b.N; n++ {
			stack.Push(item)
		}
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			stack.Pop()
		}
	})
}
