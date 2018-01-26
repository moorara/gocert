package util

import (
	"math"
)

const maxUint = ^uint(0)
const minUint = 0
const maxInt = int(maxUint >> 1)
const minInt = -maxInt - 1
const maxFloat64 = math.MaxFloat64
const minFloat64 = -1 * math.MaxFloat64

// MinInt returns minimun of int numbers
func MinInt(nums ...int) int {
	if len(nums) == 0 {
		return minInt
	}

	min := nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] < min {
			min = nums[i]
		}
	}
	return min
}

// MinFloat64 returns minimun of float64 numbers
func MinFloat64(nums ...float64) float64 {
	if len(nums) == 0 {
		return minFloat64
	}

	min := nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] < min {
			min = nums[i]
		}
	}
	return min
}

// MaxInt returns maximum of int numbers
func MaxInt(nums ...int) int {
	if len(nums) == 0 {
		return maxInt
	}

	max := nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] > max {
			max = nums[i]
		}
	}
	return max
}

// MaxFloat64 returns maximum of float64 numbers
func MaxFloat64(nums ...float64) float64 {
	if len(nums) == 0 {
		return maxFloat64
	}

	max := nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] > max {
			max = nums[i]
		}
	}
	return max
}

// IsIntIn checks if an integer is in a list of integers
func IsIntIn(n int, list ...int) bool {
	for _, i := range list {
		if i == n {
			return true
		}
	}

	return false
}
