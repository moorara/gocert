package graphviz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEdge(t *testing.T) {
	tests := []struct {
		from            string
		to              string
		edgeType        string
		edgeDir         string
		label           string
		color           string
		style           string
		arrowhead       string
		expectedDotCode string
	}{
		{
			"root", "left", EdgeTypeDirected, "",
			"", "", "", "",
			`root -> left [];`,
		},
		{
			"root", "right", EdgeTypeUndirected, "",
			"normal", "", "", "",
			`root -- right [label="normal"];`,
		},
		{
			"parent", "child", EdgeTypeDirected, EdgeDirNone,
			"red", ColorGold, StyleDashed, ArrowheadBox,
			`parent -> child [dirType=none, label="red", color=gold, style=dashed, arrowhead=box];`,
		},
		{
			"parent", "child", EdgeTypeUndirected, EdgeDirBoth,
			"black", ColorOrchid, StyleDotted, ArrowheadOpen,
			`parent -- child [dirType=both, label="black", color=orchid, style=dotted, arrowhead=open];`,
		},
	}

	for _, test := range tests {
		e := NewEdge(test.from, test.to, test.edgeType, test.edgeDir, test.label, test.color, test.style, test.arrowhead)

		assert.Equal(t, test.expectedDotCode, e.DotCode())
	}
}
