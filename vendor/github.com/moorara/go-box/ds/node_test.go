package ds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewArrayNodeWithSizeAndNext(t *testing.T) {
	tests := []struct {
		size int
		next *arrayNode
	}{
		{64, nil},
		{256, nil},
		{1024, &arrayNode{}},
		{4096, &arrayNode{}},
	}

	for _, test := range tests {
		n := newArrayNode(test.size, test.next)

		assert.Equal(t, test.next, n.next)
		assert.Equal(t, test.size, len(n.block))
	}
}
