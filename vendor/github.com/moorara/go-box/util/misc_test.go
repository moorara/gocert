package util

import (
	"testing"

	. "github.com/moorara/go-box/dt"
	"github.com/stretchr/testify/assert"
)

func TestIsSorted(t *testing.T) {
	tests := []struct {
		compare        Compare
		items          []Generic
		expectedSorted bool
	}{
		{CompareInt, []Generic{}, true},
		{CompareInt, []Generic{10, 20, 30, 50, 40}, false},
		{CompareInt, []Generic{10, 20, 30, 40, 50, 60, 70, 80, 90}, true},
		{CompareString, []Generic{"Alice", "Bob", "Dan", "Edgar", "Helen", "Karen", "Milad", "Peter", "Sam", "Wesley"}, true},
	}

	for _, test := range tests {
		assert.Equal(t, test.expectedSorted, IsSorted(test.items, test.compare))
	}
}

func TestIsStringIn(t *testing.T) {
	tests := []struct {
		str            string
		list           []string
		expectedResult bool
	}{
		{"Alice", []string{}, false},
		{"Alice", []string{"Alice"}, true},
		{"Alice", []string{"Bob", "Jackie"}, false},
		{"Jackie", []string{"Bob", "Jackie"}, true},
	}

	for _, test := range tests {
		assert.Equal(t, test.expectedResult, IsStringIn(test.str, test.list...))
	}
}
