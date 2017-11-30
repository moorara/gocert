package heap

import (
	. "github.com/moorara/go-box/dt"
)

type maxHeap struct {
	last         int
	keys         []Generic
	values       []Generic
	compareKey   Compare
	compareValue Compare
}

// NewMaxHeap creates a new max-heap (priority queue)
func NewMaxHeap(initialSize int, compareKey, compareValue Compare) Heap {
	return &maxHeap{
		last:         0,
		keys:         make([]Generic, initialSize),
		values:       make([]Generic, initialSize),
		compareKey:   compareKey,
		compareValue: compareValue,
	}
}

func (h *maxHeap) resize(newSize int) {
	newKeys := make([]Generic, newSize)
	newValues := make([]Generic, newSize)

	copy(newKeys, h.keys)
	copy(newValues, h.values)

	h.keys = newKeys
	h.values = newValues
}

func (h *maxHeap) Size() int {
	return h.last
}

func (h *maxHeap) IsEmpty() bool {
	return h.last == 0
}

func (h *maxHeap) Insert(key Generic, value Generic) {
	if h.last == len(h.keys)-1 {
		h.resize(len(h.keys) * 2)
	}

	h.last++
	var i int

	for i = h.last; true; i /= 2 {
		if i == 1 || h.compareKey(key, h.keys[i/2]) <= 0 {
			break
		}
		h.keys[i] = h.keys[i/2]
		h.values[i] = h.values[i/2]
	}

	h.keys[i] = key
	h.values[i] = value
}

func (h *maxHeap) Delete() (Generic, Generic) {
	if h.last == 0 {
		return nil, nil
	}

	maxKey := h.keys[1]
	maxValue := h.values[1]
	lastKey := h.keys[h.last]
	lastValue := h.values[h.last]

	h.last--
	var i, j int

	for i, j = 1, 2; j <= h.last; i, j = j, j*2 {
		if j < h.last && h.compareKey(h.keys[j], h.keys[j+1]) < 0 {
			j++
		}
		if h.compareKey(lastKey, h.keys[j]) >= 0 {
			break
		}
		h.keys[i] = h.keys[j]
		h.values[i] = h.values[j]
	}

	h.keys[i] = lastKey
	h.values[i] = lastValue

	if h.last < len(h.keys)/4 {
		h.resize(len(h.keys) / 2)
	}

	return maxKey, maxValue
}

func (h *maxHeap) Peek() (Generic, Generic) {
	if h.last == 0 {
		return nil, nil
	}

	return h.keys[1], h.values[1]
}

func (h *maxHeap) ContainsKey(key Generic) bool {
	for i := 1; i <= h.last; i++ {
		if h.compareKey(h.keys[i], key) == 0 {
			return true
		}
	}

	return false
}

func (h *maxHeap) ContainsValue(value Generic) bool {
	for i := 1; i <= h.last; i++ {
		if h.compareValue(h.values[i], value) == 0 {
			return true
		}
	}

	return false
}
