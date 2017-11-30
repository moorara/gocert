package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinInt(t *testing.T) {
	tests := []struct {
		nums        []int
		expectedMin int
	}{
		{[]int{}, minInt},
		{[]int{7}, 7},
		{[]int{10, 20}, 10},
		{[]int{40, 50, 20, 30, 10}, 10},
	}

	for _, test := range tests {
		assert.Equal(t, test.expectedMin, MinInt(test.nums...))
	}
}

func TestMinFloat64(t *testing.T) {
	tests := []struct {
		nums        []float64
		expectedMin float64
	}{
		{[]float64{}, minFloat64},
		{[]float64{3.14}, 3.14},
		{[]float64{.10, .20}, .10},
		{[]float64{.40, .50, .20, .30, .10}, .10},
	}

	for _, test := range tests {
		assert.Equal(t, test.expectedMin, MinFloat64(test.nums...))
	}
}

func TestMaxInt(t *testing.T) {
	tests := []struct {
		nums        []int
		expectedMax int
	}{
		{[]int{}, maxInt},
		{[]int{7}, 7},
		{[]int{10, 20}, 20},
		{[]int{40, 50, 20, 30, 10}, 50},
	}

	for _, test := range tests {
		assert.Equal(t, test.expectedMax, MaxInt(test.nums...))
	}
}

func TestMaxFloat64(t *testing.T) {
	tests := []struct {
		nums        []float64
		expectedMax float64
	}{
		{[]float64{}, maxFloat64},
		{[]float64{3.14}, 3.14},
		{[]float64{.10, .20}, .20},
		{[]float64{.40, .50, .20, .30, .10}, .50},
	}

	for _, test := range tests {
		assert.Equal(t, test.expectedMax, MaxFloat64(test.nums...))
	}
}

func TestIsIntIn(t *testing.T) {
	tests := []struct {
		num            int
		list           []int
		expectedResult bool
	}{
		{5, []int{}, false},
		{5, []int{5}, true},
		{5, []int{10, 20}, false},
		{10, []int{10, 20}, true},
	}

	for _, test := range tests {
		assert.Equal(t, test.expectedResult, IsIntIn(test.num, test.list...))
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
