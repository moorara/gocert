package ds

import (
	. "github.com/moorara/go-box/dt"
)

// Queue represents a stack data structure
type Queue interface {
	Size() int
	IsEmpty() bool
	Enqueue(Generic)
	Dequeue() Generic
	Peek() Generic
	Contains(Generic) bool
}

type arrayQueue struct {
	listSize       int
	nodeSize       int
	frontNodeIndex int
	rearNodeIndex  int
	frontNode      *arrayNode
	rearNode       *arrayNode
	compare        Compare
}

// NewQueue creates a new array-list queue
func NewQueue(nodeSize int, compare Compare) Queue {
	return &arrayQueue{
		listSize:       0,
		nodeSize:       nodeSize,
		frontNodeIndex: -1,
		rearNodeIndex:  -1,
		frontNode:      nil,
		rearNode:       nil,
		compare:        compare,
	}
}

func (q *arrayQueue) Size() int {
	return q.listSize
}

func (q *arrayQueue) IsEmpty() bool {
	return q.listSize == 0
}

func (q *arrayQueue) Enqueue(item Generic) {
	if q.frontNode == nil && q.rearNode == nil {
		q.frontNodeIndex = 0
		q.rearNodeIndex = 0
		q.frontNode = newArrayNode(q.nodeSize, nil)
		q.rearNode = q.frontNode
	}

	q.listSize++
	q.rearNode.block[q.rearNodeIndex] = item
	q.rearNodeIndex++

	if q.rearNodeIndex == q.nodeSize {
		q.rearNodeIndex = 0
		q.rearNode.next = newArrayNode(q.nodeSize, nil)
		q.rearNode = q.rearNode.next
	}
}

func (q *arrayQueue) Dequeue() Generic {
	if q.listSize == 0 {
		return nil
	}

	q.listSize--
	item := q.frontNode.block[q.frontNodeIndex]
	q.frontNodeIndex++

	if q.frontNodeIndex == q.nodeSize {
		q.frontNodeIndex = 0
		q.frontNode = q.frontNode.next
	}

	return item
}

func (q *arrayQueue) Peek() Generic {
	if q.listSize == 0 {
		return nil
	}

	return q.frontNode.block[q.frontNodeIndex]
}

func (q *arrayQueue) Contains(item Generic) bool {
	n := q.frontNode
	i := q.frontNodeIndex

	for n != nil && (n != q.rearNode || i <= q.rearNodeIndex) {
		if q.compare(n.block[i], item) == 0 {
			return true
		}

		i++
		if i == q.nodeSize {
			n = n.next
			i = 0
		}
	}

	return false
}
