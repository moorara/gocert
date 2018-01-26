package util

import (
	"testing"

	. "github.com/moorara/goto/dt"
	"github.com/stretchr/testify/assert"
)

func TestShuffle(t *testing.T) {
	tests := []struct {
		items []Generic
	}{
		{[]Generic{10, 20, 30, 40, 50, 60, 70, 80, 90}},
		{[]Generic{"Alice", "Bob", "Dan", "Edgar", "Helen", "Karen", "Milad", "Peter", "Sam", "Wesley"}},
	}

	SeedWithNow()

	for _, test := range tests {
		orig := make([]Generic, len(test.items))
		copy(orig, test.items)
		Shuffle(test.items)

		assert.NotEqual(t, orig, test.items)
	}
}

func TestGenerateInt(t *testing.T) {
	tests := []struct {
		min int
		max int
	}{
		{0, 0},
		{1, 1},
		{0, 1000},
		{100, 100000},
	}

	SeedWithNow()

	for _, test := range tests {
		n := GenerateInt(test.min, test.max)

		assert.True(t, test.min <= n && n <= test.max)
	}
}

func TestGenerateString(t *testing.T) {
	tests := []struct {
		minLen int
		maxLen int
	}{
		{0, 0},
		{1, 1},
		{10, 100},
		{100, 1000},
	}

	SeedWithNow()

	for _, test := range tests {
		str := GenerateString(test.minLen, test.maxLen)

		assert.True(t, test.minLen <= len(str) && len(str) <= test.maxLen)
	}
}

func TestGenerateIntSlice(t *testing.T) {
	tests := []struct {
		size int
		min  int
		max  int
	}{
		{0, 0, 0},
		{1, 1, 1},
		{10, 0, 100},
		{100, 100, 1000},
	}

	SeedWithNow()

	for _, test := range tests {
		items := GenerateIntSlice(test.size, test.min, test.max)
		for _, item := range items {
			if CompareInt(item, test.min) < 0 || CompareInt(item, test.max) > 0 {
				t.Errorf("%d is not between %d and %d.", item, test.min, test.max)
			}
		}
	}
}

func TestGenerateStringSlice(t *testing.T) {
	tests := []struct {
		size   int
		minLen int
		maxLen int
	}{
		{0, 0, 0},
		{1, 1, 1},
		{10, 1, 10},
		{100, 10, 100},
	}

	SeedWithNow()

	for _, test := range tests {
		items := GenerateStringSlice(test.size, test.minLen, test.maxLen)
		for _, item := range items {
			if len(item.(string)) < test.minLen || len(item.(string)) > test.maxLen {
				t.Errorf("%s length is not between %d and %d.", item, test.minLen, test.maxLen)
			}
		}
	}
}
