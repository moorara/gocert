package st

import (
	. "github.com/moorara/go-box/dt"
)

const (
	// TraversePreOrder represents pre-order traversal order
	TraversePreOrder = 0
	// TraverseInOrder represents in-order traversal order
	TraverseInOrder = 1
	// TraversePostOrder represents post-order traversal order
	TraversePostOrder = 2
)

type (
	// VisitFunc represents the function for visting a key-value
	VisitFunc func(Generic, Generic) bool

	// KeyValue represents a key-value pair
	KeyValue struct {
		key   Generic
		value Generic
	}

	// SymbolTable represents an unordered symbol table (key-value collection)
	SymbolTable interface {
		verify() bool
		Size() int
		Height() int
		IsEmpty() bool
		Put(Generic, Generic)
		Get(Generic) (Generic, bool)
		Delete(Generic) (Generic, bool)
		KeyValues() []KeyValue
	}

	// OrderedSymbolTable represents an ordered symbol table (key-value collection)
	OrderedSymbolTable interface {
		SymbolTable
		Min() (Generic, Generic)
		Max() (Generic, Generic)
		Floor(Generic) (Generic, Generic)
		Ceiling(Generic) (Generic, Generic)
		Rank(Generic) int
		Select(int) (Generic, Generic)
		DeleteMin() (Generic, Generic)
		DeleteMax() (Generic, Generic)
		RangeSize(Generic, Generic) int
		Range(Generic, Generic) []KeyValue
		Traverse(int, VisitFunc)
		Graphviz() string
	}
)
