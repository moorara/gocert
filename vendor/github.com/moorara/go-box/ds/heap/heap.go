package heap

import (
	. "github.com/moorara/go-box/dt"
)

// Heap represents a heap (priority queue) data structure
type Heap interface {
	Size() int
	IsEmpty() bool
	Insert(Generic, Generic)
	Delete() (Generic, Generic)
	Peek() (Generic, Generic)
	ContainsKey(Generic) bool
	ContainsValue(Generic) bool
}
