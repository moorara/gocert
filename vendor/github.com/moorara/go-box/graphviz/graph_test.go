package graphviz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraph(t *testing.T) {
	tests := []struct {
		strict          bool
		diagraph        bool
		name            string
		rankdir         string
		nodeColor       string
		nodeStyle       string
		nodeShape       string
		nodes           []Node
		edges           []Edge
		subgraphs       []Subgraph
		expectedDotCode string
	}{
		{
			false, false, "",
			"",
			"", "", "",
			[]Node{
				Node{Name: "a0"},
				Node{Name: "a1"},
			},
			[]Edge{
				Edge{From: "a0", To: "a1", EdgeType: EdgeTypeUndirected},
			},
			[]Subgraph{},
			`graph {
  node [];

  a0 [];
  a1 [];

  a0 -- a1 [];
}`,
		},
		{
			true, false, "G",
			"",
			"", "", "",
			[]Node{
				Node{Name: "b0", Label: "B0"},
				Node{Name: "b1", Label: "B3"},
				Node{Name: "b2", Label: "B2"},
			},
			[]Edge{
				Edge{From: "b0", To: "b1", EdgeType: EdgeTypeUndirected, Color: ColorRed},
				Edge{From: "b0", To: "b2", EdgeType: EdgeTypeUndirected, Color: ColorBlack},
			},
			[]Subgraph{},
			`strict graph G {
  node [];

  b0 [label="B0"];
  b1 [label="B3"];
  b2 [label="B2"];

  b0 -- b1 [color=red];
  b0 -- b2 [color=black];
}`,
		},
		{
			false, true, "",
			"",
			ColorLimeGreen, "", "",
			[]Node{
				Node{Name: "c0", Label: "C0", Shape: ShapePlain},
				Node{Name: "c1", Label: "C1", Shape: ShapePlain},
				Node{Name: "c2", Label: "C2", Shape: ShapePlain},
				Node{Name: "c3", Label: "C3", Shape: ShapePlain},
			},
			[]Edge{
				Edge{From: "c0", To: "c1", EdgeType: EdgeTypeDirected, EdgeDir: EdgeDirBoth, Arrowhead: ArrowheadOpen},
				Edge{From: "c0", To: "c2", EdgeType: EdgeTypeDirected, EdgeDir: EdgeDirBoth, Arrowhead: ArrowheadOpen},
				Edge{From: "c2", To: "c3", EdgeType: EdgeTypeDirected, EdgeDir: EdgeDirBoth, Arrowhead: ArrowheadOpen},
			},
			[]Subgraph{
				Subgraph{Name: "", Label: "Thread", Rank: RankSame},
			},
			`digraph {
  node [color=limegreen];

  subgraph {
    label="Thread";
    rank=same;
    node [];
  }

  c0 [label="C0", shape=plain];
  c1 [label="C1", shape=plain];
  c2 [label="C2", shape=plain];
  c3 [label="C3", shape=plain];

  c0 -> c1 [dirType=both, arrowhead=open];
  c0 -> c2 [dirType=both, arrowhead=open];
  c2 -> c3 [dirType=both, arrowhead=open];
}`,
		},
		{
			true, true, "DG",
			RankdirLR,
			ColorSteelBlue, StyleFilled, ShapeMrecord,
			[]Node{
				Node{Name: "start", Label: "Start", Color: ColorBlue, Shape: ShapeBox},
				Node{Name: "end", Label: "End", Color: ColorBlue, Shape: ShapeBox},
			},
			[]Edge{
				Edge{From: "start", To: "e0", EdgeType: EdgeTypeDirected, Label: "Start", Color: ColorRed, Style: StyleSolid},
				Edge{From: "start", To: "f0", EdgeType: EdgeTypeDirected, Label: "Start", Color: ColorRed, Style: StyleSolid},
				Edge{From: "e1", To: "end", EdgeType: EdgeTypeDirected, Label: "End", Color: ColorRed, Style: StyleSolid},
				Edge{From: "f1", To: "end", EdgeType: EdgeTypeDirected, Label: "End", Color: ColorRed, Style: StyleSolid},
			},
			[]Subgraph{
				Subgraph{
					Name: "cluster0", Label: "Process 0", Color: ColorGray, Style: StyleFilled,
					Nodes: []Node{
						Node{Name: "e0"},
						Node{Name: "e1"},
					},
					Edges: []Edge{
						Edge{From: "e0", To: "e1", EdgeType: EdgeTypeDirected, Style: StyleDashed},
					},
				},
				Subgraph{
					Name: "cluster1", Label: "Process 1", Color: ColorGray, Style: StyleFilled,
					Nodes: []Node{
						Node{Name: "f0"},
						Node{Name: "f1"},
					},
					Edges: []Edge{
						Edge{From: "f0", To: "f1", EdgeType: EdgeTypeDirected, Style: StyleDashed},
					},
				},
			},
			`strict digraph DG {
  rankdir=LR;
  node [color=steelblue, style=filled, shape=Mrecord];

  subgraph cluster0 {
    label="Process 0";
    color=gray;
    style=filled;
    node [];

    e0 [];
    e1 [];

    e0 -> e1 [style=dashed];
  }

  subgraph cluster1 {
    label="Process 1";
    color=gray;
    style=filled;
    node [];

    f0 [];
    f1 [];

    f0 -> f1 [style=dashed];
  }

  start [label="Start", color=blue, shape=box];
  end [label="End", color=blue, shape=box];

  start -> e0 [label="Start", color=red, style=solid];
  start -> f0 [label="Start", color=red, style=solid];
  e1 -> end [label="End", color=red, style=solid];
  f1 -> end [label="End", color=red, style=solid];
}`,
		},
	}

	for _, test := range tests {
		g := NewGraph(test.strict, test.diagraph, test.name, test.rankdir, test.nodeColor, test.nodeStyle, test.nodeShape)
		g.AddNode(test.nodes...)
		g.AddEdge(test.edges...)
		g.AddSubgraph(test.subgraphs...)

		assert.Equal(t, test.expectedDotCode, g.DotCode())
	}
}
