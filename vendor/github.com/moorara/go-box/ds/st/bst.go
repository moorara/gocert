/*
 * A Binary Search Tree (BST) is a binary tree in symmetric order.
 * Every node's key is:
 *  - larger than all keys in its left sub-tree
 *  - smaller than all keys in its right sub-tree
 */

package st

import (
	"fmt"

	. "github.com/moorara/go-box/dt"
	"github.com/moorara/go-box/graphviz"
	"github.com/moorara/go-box/util"
)

type bstNode struct {
	key   Generic
	value Generic
	left  *bstNode
	right *bstNode
	size  int
}

type bst struct {
	root       *bstNode
	compareKey Compare
}

// NewBST creates a new Binary Search Tree
func NewBST(compareKey Compare) OrderedSymbolTable {
	return &bst{
		root:       nil,
		compareKey: compareKey,
	}
}

func (t *bst) isBST(n *bstNode, min, max Generic) bool {
	if n == nil {
		return true
	}

	if (min != nil && t.compareKey(n.key, min) <= 0) ||
		(max != nil && t.compareKey(n.key, max) >= 0) {
		return false
	}

	return t.isBST(n.left, min, n.key) && t.isBST(n.right, n.key, max)
}

func (t *bst) isSizeOK(n *bstNode) bool {
	if n == nil {
		return true
	}

	if n.size != 1+t.size(n.left)+t.size(n.right) {
		return false
	}

	return t.isSizeOK(n.left) && t.isSizeOK(n.right)
}

func (t *bst) verify() bool {
	return t.isBST(t.root, nil, nil) &&
		t.isSizeOK(t.root)
}

func (t *bst) size(n *bstNode) int {
	if n == nil {
		return 0
	}

	return n.size
}

func (t *bst) height(n *bstNode) int {
	if n == nil {
		return 0
	}

	return 1 + util.MaxInt(t.height(n.left), t.height(n.right))
}

func (t *bst) Size() int {
	return t.size(t.root)
}

func (t *bst) Height() int {
	return t.height(t.root)
}

func (t *bst) IsEmpty() bool {
	return t.root == nil
}

func (t *bst) _put(n *bstNode, key, value Generic) *bstNode {
	if n == nil {
		return &bstNode{
			key:   key,
			value: value,
			size:  1,
		}
	}

	cmp := t.compareKey(key, n.key)
	switch {
	case cmp < 0:
		n.left = t._put(n.left, key, value)
	case cmp > 0:
		n.right = t._put(n.right, key, value)
	default:
		n.value = value
	}

	n.size = 1 + t.size(n.left) + t.size(n.right)
	return n
}

func (t *bst) Put(key, value Generic) {
	if key == nil {
		return
	}

	t.root = t._put(t.root, key, value)
}

func (t *bst) _get(n *bstNode, key Generic) (Generic, bool) {
	if n == nil || key == nil {
		return nil, false
	}

	cmp := t.compareKey(key, n.key)
	switch {
	case cmp < 0:
		return t._get(n.left, key)
	case cmp > 0:
		return t._get(n.right, key)
	default:
		return n.value, true
	}
}

func (t *bst) Get(key Generic) (Generic, bool) {
	return t._get(t.root, key)
}

func (t *bst) _delete(n *bstNode, key Generic) (*bstNode, Generic, bool) {
	if n == nil || key == nil {
		return n, nil, false
	}

	var ok bool
	var value Generic

	cmp := t.compareKey(key, n.key)
	if cmp < 0 {
		n.left, value, ok = t._delete(n.left, key)
	} else if cmp > 0 {
		n.right, value, ok = t._delete(n.right, key)
	} else {
		ok = true
		value = n.value

		if n.left == nil {
			return n.right, value, ok
		} else if n.right == nil {
			return n.left, value, ok
		} else {
			m := n
			n = t._min(m.right)
			n.right, _ = t._deleteMin(m.right)
			n.left = m.left
		}
	}

	n.size = 1 + t.size(n.left) + t.size(n.right)
	return n, value, ok
}

func (t *bst) Delete(key Generic) (value Generic, ok bool) {
	t.root, value, ok = t._delete(t.root, key)
	return value, ok
}

func (t *bst) KeyValues() []KeyValue {
	i := 0
	kvs := make([]KeyValue, t.Size())

	t._traverse(t.root, TraverseInOrder, func(n *bstNode) bool {
		kvs[i] = KeyValue{n.key, n.value}
		i++
		return true
	})
	return kvs
}

func (t *bst) _min(n *bstNode) *bstNode {
	if n.left == nil {
		return n
	}
	return t._min(n.left)
}

func (t *bst) Min() (Generic, Generic) {
	if t.root == nil {
		return nil, nil
	}

	n := t._min(t.root)
	return n.key, n.value
}

func (t *bst) _max(n *bstNode) *bstNode {
	if n.right == nil {
		return n
	}
	return t._max(n.right)
}

func (t *bst) Max() (Generic, Generic) {
	if t.root == nil {
		return nil, nil
	}

	n := t._max(t.root)
	return n.key, n.value
}

func (t *bst) _floor(n *bstNode, key Generic) *bstNode {
	if n == nil || key == nil {
		return nil
	}

	cmp := t.compareKey(key, n.key)
	if cmp == 0 {
		return n
	} else if cmp < 0 {
		return t._floor(n.left, key)
	}

	m := t._floor(n.right, key)
	if m != nil {
		return m
	}
	return n
}

func (t *bst) Floor(key Generic) (Generic, Generic) {
	n := t._floor(t.root, key)
	if n == nil {
		return nil, nil
	}
	return n.key, n.value
}

func (t *bst) _ceiling(n *bstNode, key Generic) *bstNode {
	if n == nil || key == nil {
		return nil
	}

	cmp := t.compareKey(key, n.key)
	if cmp == 0 {
		return n
	} else if cmp > 0 {
		return t._ceiling(n.right, key)
	}

	m := t._ceiling(n.left, key)
	if m != nil {
		return m
	}
	return n
}

func (t *bst) Ceiling(key Generic) (Generic, Generic) {
	n := t._ceiling(t.root, key)
	if n == nil {
		return nil, nil
	}
	return n.key, n.value
}

func (t *bst) _rank(n *bstNode, key Generic) int {
	if n == nil {
		return 0
	}

	cmp := t.compareKey(key, n.key)
	switch {
	case cmp < 0:
		return t._rank(n.left, key)
	case cmp > 0:
		return 1 + t.size(n.left) + t._rank(n.right, key)
	default:
		return t.size(n.left)
	}
}

func (t *bst) Rank(key Generic) int {
	if key == nil {
		return -1
	}

	return t._rank(t.root, key)
}

func (t *bst) _select(n *bstNode, rank int) *bstNode {
	if n == nil {
		return nil
	}

	s := t.size(n.left)
	switch {
	case rank < s:
		return t._select(n.left, rank)
	case rank > s:
		return t._select(n.right, rank-s-1)
	default:
		return n
	}
}

func (t *bst) Select(rank int) (Generic, Generic) {
	if rank < 0 || rank >= t.Size() {
		return nil, nil
	}

	n := t._select(t.root, rank)
	return n.key, n.value
}

func (t *bst) _deleteMin(n *bstNode) (*bstNode, *bstNode) {
	if n.left == nil {
		return n.right, n
	}

	var min *bstNode
	n.left, min = t._deleteMin(n.left)
	n.size = 1 + t.size(n.left) + t.size(n.right)
	return n, min
}

func (t *bst) DeleteMin() (Generic, Generic) {
	if t.root == nil {
		return nil, nil
	}

	var min *bstNode
	t.root, min = t._deleteMin(t.root)
	return min.key, min.value
}

func (t *bst) _deleteMax(n *bstNode) (*bstNode, *bstNode) {
	if n.right == nil {
		return n.left, n
	}

	var max *bstNode
	n.right, max = t._deleteMax(n.right)
	n.size = 1 + t.size(n.left) + t.size(n.right)
	return n, max
}

func (t *bst) DeleteMax() (Generic, Generic) {
	if t.root == nil {
		return nil, nil
	}

	var max *bstNode
	t.root, max = t._deleteMax(t.root)
	return max.key, max.value
}

func (t *bst) RangeSize(lo, hi Generic) int {
	if lo == nil || hi == nil {
		return -1
	}

	if t.compareKey(lo, hi) > 0 {
		return 0
	} else if _, found := t.Get(hi); found {
		return 1 + t.Rank(hi) - t.Rank(lo)
	} else {
		return t.Rank(hi) - t.Rank(lo)
	}
}

func (t *bst) _range(n *bstNode, kvs *[]KeyValue, lo, hi Generic) int {
	if n == nil {
		return 0
	}

	len := 0
	cmpLo := t.compareKey(lo, n.key)
	cmpHi := t.compareKey(hi, n.key)

	if cmpLo < 0 {
		len += t._range(n.left, kvs, lo, hi)
	}
	if cmpLo <= 0 && cmpHi >= 0 {
		*kvs = append(*kvs, KeyValue{n.key, n.value})
		len++
	}
	if cmpHi > 0 {
		len += t._range(n.right, kvs, lo, hi)
	}

	return len
}

func (t *bst) Range(lo, hi Generic) []KeyValue {
	if lo == nil || hi == nil {
		return nil
	}

	kvs := make([]KeyValue, 0)
	len := t._range(t.root, &kvs, lo, hi)
	return kvs[0:len]
}

func (t *bst) _traverse(n *bstNode, order int, visit func(*bstNode) bool) bool {
	if n == nil {
		return true
	}

	switch order {
	case TraversePreOrder:
		return visit(n) &&
			t._traverse(n.left, order, visit) &&
			t._traverse(n.right, order, visit)
	case TraverseInOrder:
		return t._traverse(n.left, order, visit) &&
			visit(n) &&
			t._traverse(n.right, order, visit)
	case TraversePostOrder:
		return t._traverse(n.left, order, visit) &&
			t._traverse(n.right, order, visit) &&
			visit(n)
	default:
		return false
	}
}

func (t *bst) Traverse(order int, visit VisitFunc) {
	if !util.IsIntIn(order, TraversePreOrder, TraverseInOrder, TraversePostOrder) {
		return
	}

	t._traverse(t.root, order, func(n *bstNode) bool {
		return visit(n.key, n.value)
	})
}

func (t *bst) Graphviz() string {
	var parent, left, right, label string
	graph := graphviz.NewGraph(true, true, "BST", "", "", "", graphviz.ShapeOval)

	t._traverse(t.root, TraversePreOrder, func(n *bstNode) bool {
		parent = fmt.Sprintf("%v", n.key)
		label = fmt.Sprintf("%v,%v", n.key, n.value)
		graph.AddNode(graphviz.NewNode(parent, "", label, "", "", "", "", ""))
		if n.left != nil {
			left = fmt.Sprintf("%v", n.left.key)
			graph.AddEdge(graphviz.NewEdge(parent, left, graphviz.EdgeTypeDirected, "", "", "", "", ""))
		}
		if n.right != nil {
			right = fmt.Sprintf("%v", n.right.key)
			graph.AddEdge(graphviz.NewEdge(parent, right, graphviz.EdgeTypeDirected, "", "", "", "", ""))
		}
		return true
	})

	return graph.DotCode()
}
