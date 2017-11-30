package st

import (
	"testing"
)

func getRedBlackTests() []orderedSymbolTableTest {
	tests := getOrderedSymbolTableTests()

	tests[0].SymbolTable = "LLRB Tree"
	tests[0].expectedHeight = 0
	tests[0].expectedPreOrderTraverse = nil
	tests[0].expectedInOrderTraverse = nil
	tests[0].expectedPostOrderTraverse = nil
	tests[0].expectedDotCode = `strict digraph RedBlack {
  node [style=filled, shape=oval];
}`

	tests[1].SymbolTable = "LLRB Tree"
	tests[1].expectedHeight = 2
	tests[1].expectedPreOrderTraverse = []KeyValue{{"B", 2}, {"A", 1}, {"C", 3}}
	tests[1].expectedInOrderTraverse = []KeyValue{{"A", 1}, {"B", 2}, {"C", 3}}
	tests[1].expectedPostOrderTraverse = []KeyValue{{"A", 1}, {"C", 3}, {"B", 2}}
	tests[1].expectedDotCode = `strict digraph RedBlack {
  node [style=filled, shape=oval];

  B [label="B,2", color=black, fontcolor=white];
  A [label="A,1", color=black, fontcolor=white];
  C [label="C,3", color=black, fontcolor=white];

  B -> A [color=black];
  B -> C [];
}`

	tests[2].SymbolTable = "LLRB Tree"
	tests[2].expectedHeight = 3
	tests[2].expectedPreOrderTraverse = []KeyValue{{"D", 4}, {"B", 2}, {"A", 1}, {"C", 3}, {"E", 5}}
	tests[2].expectedInOrderTraverse = []KeyValue{{"A", 1}, {"B", 2}, {"C", 3}, {"D", 4}, {"E", 5}}
	tests[2].expectedPostOrderTraverse = []KeyValue{{"A", 1}, {"C", 3}, {"B", 2}, {"E", 5}, {"D", 4}}
	tests[2].expectedDotCode = `strict digraph RedBlack {
  node [style=filled, shape=oval];

  D [label="D,4", color=black, fontcolor=white];
  B [label="B,2", color=red, fontcolor=white];
  A [label="A,1", color=black, fontcolor=white];
  C [label="C,3", color=black, fontcolor=white];
  E [label="E,5", color=black, fontcolor=white];

  D -> B [color=red];
  D -> E [];
  B -> A [color=black];
  B -> C [];
}`

	tests[3].SymbolTable = "LLRB Tree"
	tests[3].expectedHeight = 3
	tests[3].expectedPreOrderTraverse = []KeyValue{{"J", 10}, {"D", 4}, {"A", 1}, {"G", 7}, {"P", 16}, {"M", 13}, {"S", 19}}
	tests[3].expectedInOrderTraverse = []KeyValue{{"A", 1}, {"D", 4}, {"G", 7}, {"J", 10}, {"M", 13}, {"P", 16}, {"S", 19}}
	tests[3].expectedPostOrderTraverse = []KeyValue{{"A", 1}, {"G", 7}, {"D", 4}, {"M", 13}, {"S", 19}, {"P", 16}, {"J", 10}}
	tests[3].expectedDotCode = `strict digraph RedBlack {
  node [style=filled, shape=oval];

  J [label="J,10", color=black, fontcolor=white];
  D [label="D,4", color=black, fontcolor=white];
  A [label="A,1", color=black, fontcolor=white];
  G [label="G,7", color=black, fontcolor=white];
  P [label="P,16", color=black, fontcolor=white];
  M [label="M,13", color=black, fontcolor=white];
  S [label="S,19", color=black, fontcolor=white];

  J -> D [color=black];
  J -> P [];
  D -> A [color=black];
  D -> G [];
  P -> M [color=black];
  P -> S [];
}`

	return tests
}

func TestRedBlack(t *testing.T) {
	tests := getRedBlackTests()

	for _, test := range tests {
		rbt := NewRedBlack(test.compareKey)
		runOrderedSymbolTableTest(t, rbt, test)
	}
}
