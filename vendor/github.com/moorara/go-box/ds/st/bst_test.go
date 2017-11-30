package st

import (
	"testing"
)

func getBSTTests() []orderedSymbolTableTest {
	tests := getOrderedSymbolTableTests()

	tests[0].SymbolTable = "BST"
	tests[0].expectedHeight = 0
	tests[0].expectedPreOrderTraverse = nil
	tests[0].expectedInOrderTraverse = nil
	tests[0].expectedPostOrderTraverse = nil
	tests[0].expectedDotCode = `strict digraph BST {
  node [shape=oval];
}`

	tests[1].SymbolTable = "BST"
	tests[1].expectedHeight = 2
	tests[1].expectedPreOrderTraverse = []KeyValue{{"B", 2}, {"A", 1}, {"C", 3}}
	tests[1].expectedInOrderTraverse = []KeyValue{{"A", 1}, {"B", 2}, {"C", 3}}
	tests[1].expectedPostOrderTraverse = []KeyValue{{"A", 1}, {"C", 3}, {"B", 2}}
	tests[1].expectedDotCode = `strict digraph BST {
  node [shape=oval];

  B [label="B,2"];
  A [label="A,1"];
  C [label="C,3"];

  B -> A [];
  B -> C [];
}`

	tests[2].SymbolTable = "BST"
	tests[2].expectedHeight = 4
	tests[2].expectedPreOrderTraverse = []KeyValue{{"B", 2}, {"A", 1}, {"C", 3}, {"D", 4}, {"E", 5}}
	tests[2].expectedInOrderTraverse = []KeyValue{{"A", 1}, {"B", 2}, {"C", 3}, {"D", 4}, {"E", 5}}
	tests[2].expectedPostOrderTraverse = []KeyValue{{"A", 1}, {"E", 5}, {"D", 4}, {"C", 3}, {"B", 2}}
	tests[2].expectedDotCode = `strict digraph BST {
  node [shape=oval];

  B [label="B,2"];
  A [label="A,1"];
  C [label="C,3"];
  D [label="D,4"];
  E [label="E,5"];

  B -> A [];
  B -> C [];
  C -> D [];
  D -> E [];
}`

	tests[3].SymbolTable = "BST"
	tests[3].expectedHeight = 4
	tests[3].expectedPreOrderTraverse = []KeyValue{{"J", 10}, {"D", 4}, {"A", 1}, {"G", 7}, {"P", 16}, {"M", 13}, {"S", 19}}
	tests[3].expectedInOrderTraverse = []KeyValue{{"A", 1}, {"D", 4}, {"G", 7}, {"J", 10}, {"M", 13}, {"P", 16}, {"S", 19}}
	tests[3].expectedPostOrderTraverse = []KeyValue{{"A", 1}, {"G", 7}, {"D", 4}, {"M", 13}, {"S", 19}, {"P", 16}, {"J", 10}}
	tests[3].expectedDotCode = `strict digraph BST {
  node [shape=oval];

  J [label="J,10"];
  D [label="D,4"];
  A [label="A,1"];
  G [label="G,7"];
  P [label="P,16"];
  M [label="M,13"];
  S [label="S,19"];

  J -> D [];
  J -> P [];
  D -> A [];
  D -> G [];
  P -> M [];
  P -> S [];
}`

	return tests
}

func TestBST(t *testing.T) {
	tests := getBSTTests()

	for _, test := range tests {
		bst := NewBST(test.compareKey)
		runOrderedSymbolTableTest(t, bst, test)
	}
}
