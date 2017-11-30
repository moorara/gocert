package graphviz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubgraph(t *testing.T) {
	tests := []struct {
		name            string
		label           string
		color           string
		style           string
		rank            string
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
			"parent",
			"", "", "",
			"", "",
			"", "", "",
			[]Node{
				Node{Name: "a0"},
				Node{Name: "a1"},
			},
			[]Edge{
				Edge{From: "a0", To: "a1", EdgeType: EdgeTypeUndirected},
			},
			[]Subgraph{},
			`subgraph parent {
  node [];

  a0 [];
  a1 [];

  a0 -- a1 [];
}`,
		},
		{
			"child",
			"Child", "", "",
			"", "",
			"", "", "",
			[]Node{
				Node{Name: "b0", Label: "B0"},
				Node{Name: "b1", Label: "B3"},
				Node{Name: "b2", Label: "B2"},
			},
			[]Edge{
				Edge{From: "b0", To: "b1", EdgeType: EdgeTypeDirected, Color: ColorRed},
				Edge{From: "b0", To: "b2", EdgeType: EdgeTypeDirected, Color: ColorBlack},
			},
			[]Subgraph{},
			`subgraph child {
  label="Child";
  node [];

  b0 [label="B0"];
  b1 [label="B3"];
  b2 [label="B2"];

  b0 -> b1 [color=red];
  b0 -> b2 [color=black];
}`,
		},
		{
			"cluster0",
			"Left", ColorPink, "",
			"", RankdirLR,
			ColorRoyalBlue, "", "",
			[]Node{
				Node{Name: "c0", Label: "C0", Shape: ShapeBox},
				Node{Name: "c1", Label: "C1", Shape: ShapeBox},
				Node{Name: "c2", Label: "C2", Shape: ShapeBox},
				Node{Name: "c3", Label: "C3", Shape: ShapeBox},
			},
			[]Edge{
				Edge{From: "c0", To: "c1", EdgeType: EdgeTypeUndirected, EdgeDir: EdgeDirBoth, Arrowhead: ArrowheadDot},
				Edge{From: "c0", To: "c2", EdgeType: EdgeTypeUndirected, EdgeDir: EdgeDirBoth, Arrowhead: ArrowheadDot},
				Edge{From: "c1", To: "c3", EdgeType: EdgeTypeUndirected, EdgeDir: EdgeDirBoth, Arrowhead: ArrowheadDot},
			},
			[]Subgraph{
				Subgraph{Name: "thread", Label: "Thread"},
			},
			`subgraph cluster0 {
  label="Left";
  color=pink;
  rankdir=LR;
  node [color=royalblue];

  subgraph thread {
    label="Thread";
    node [];
  }

  c0 [label="C0", shape=box];
  c1 [label="C1", shape=box];
  c2 [label="C2", shape=box];
  c3 [label="C3", shape=box];

  c0 -- c1 [dirType=both, arrowhead=dot];
  c0 -- c2 [dirType=both, arrowhead=dot];
  c1 -- c3 [dirType=both, arrowhead=dot];
}`,
		},
		{
			"cluster1",
			"Right", ColorGreen, StyleDotted,
			RankSame, "",
			ColorSeaGreen, StyleFilled, ShapeRecord,
			[]Node{
				Node{Name: "d0", Label: "D0", Color: ColorTan, Shape: ShapeOval},
				Node{Name: "d1", Label: "D1", Color: ColorTan, Shape: ShapeOval},
			},
			[]Edge{
				Edge{From: "d0", To: "e0", EdgeType: EdgeTypeDirected, Label: "d0e0", Color: ColorGray, Style: StyleDashed},
				Edge{From: "d0", To: "f0", EdgeType: EdgeTypeDirected, Label: "d0f0", Color: ColorGray, Style: StyleDashed},
				Edge{From: "e1", To: "d1", EdgeType: EdgeTypeDirected, Label: "e1d1", Color: ColorGray, Style: StyleDashed},
				Edge{From: "f1", To: "d1", EdgeType: EdgeTypeDirected, Label: "f1d1", Color: ColorGray, Style: StyleDashed},
			},
			[]Subgraph{
				Subgraph{
					Name: "process0", Label: "Process 0", Color: ColorGray, Style: StyleFilled,
					Nodes: []Node{
						Node{Name: "e0"},
						Node{Name: "e1"},
					},
					Edges: []Edge{
						Edge{From: "e0", To: "e1", EdgeType: EdgeTypeDirected},
					},
				},
				Subgraph{
					Name: "process1", Label: "Process 1", Color: ColorGray, Style: StyleFilled,
					Nodes: []Node{
						Node{Name: "f0"},
						Node{Name: "f1"},
					},
					Edges: []Edge{
						Edge{From: "f0", To: "f1", EdgeType: EdgeTypeDirected},
					},
				},
			},
			`subgraph cluster1 {
  label="Right";
  color=green;
  style=dotted;
  rank=same;
  node [color=seagreen, style=filled, shape=record];

  subgraph process0 {
    label="Process 0";
    color=gray;
    style=filled;
    node [];

    e0 [];
    e1 [];

    e0 -> e1 [];
  }

  subgraph process1 {
    label="Process 1";
    color=gray;
    style=filled;
    node [];

    f0 [];
    f1 [];

    f0 -> f1 [];
  }

  d0 [label="D0", color=tan, shape=oval];
  d1 [label="D1", color=tan, shape=oval];

  d0 -> e0 [label="d0e0", color=gray, style=dashed];
  d0 -> f0 [label="d0f0", color=gray, style=dashed];
  e1 -> d1 [label="e1d1", color=gray, style=dashed];
  f1 -> d1 [label="f1d1", color=gray, style=dashed];
}`,
		},
	}

	for _, test := range tests {
		sg := NewSubgraph(test.name, test.label, test.color, test.style, test.rank, test.rankdir, test.nodeColor, test.nodeStyle, test.nodeShape)
		sg.AddNode(test.nodes...)
		sg.AddEdge(test.edges...)
		sg.AddSubgraph(test.subgraphs...)

		assert.Equal(t, test.expectedDotCode, sg.DotCode(0))
	}
}
