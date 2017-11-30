package graphviz

import (
	"bytes"
)

// Subgraph represents a subgraph for visualization
type Subgraph struct {
	Name      string
	Label     string
	Color     string
	Style     string
	Rank      string
	Rankdir   string
	NodeColor string
	NodeStyle string
	NodeShape string
	Nodes     []Node
	Edges     []Edge
	Subgraphs []Subgraph
}

// NewSubgraph creates a new subgraph for visualization
func NewSubgraph(name, label, color, style, rank, rankdir, nodeColor, nodeStyle, nodeShape string) Subgraph {
	return Subgraph{
		Name:      name,
		Label:     label,
		Color:     color,
		Style:     style,
		Rank:      rank,
		Rankdir:   rankdir,
		NodeColor: nodeColor,
		NodeStyle: nodeStyle,
		NodeShape: nodeShape,
		Nodes:     make([]Node, 0),
		Edges:     make([]Edge, 0),
		Subgraphs: make([]Subgraph, 0),
	}
}

// AddNode adds a new node to graph for visualization
func (s *Subgraph) AddNode(nodes ...Node) {
	for _, n := range nodes {
		s.Nodes = append(s.Nodes, n)
	}
}

// AddEdge adds a new edge to graph for visualization
func (s *Subgraph) AddEdge(edges ...Edge) {
	for _, e := range edges {
		s.Edges = append(s.Edges, e)
	}
}

// AddSubgraph adds a new subgraph to graph for visualization
func (s *Subgraph) AddSubgraph(subgraphs ...Subgraph) {
	for _, sg := range subgraphs {
		s.Subgraphs = append(s.Subgraphs, sg)
	}
}

// DotCode generates Graph dot language code for subgraph
func (s *Subgraph) DotCode(indent int) string {
	first := true
	buf := new(bytes.Buffer)

	addIndent(buf, indent)
	buf.WriteString("subgraph ")
	if s.Name != "" {
		buf.WriteString(s.Name)
		buf.WriteString(" ")
	}
	buf.WriteString("{\n")

	first = addAttr(buf, first, indent+2, "label", `"`+s.Label+`"`)
	first = addAttr(buf, first, indent+2, "color", s.Color)
	first = addAttr(buf, first, indent+2, "style", s.Style)
	first = addAttr(buf, first, indent+2, "rank", s.Rank)
	first = addAttr(buf, first, indent+2, "rankdir", s.Rankdir)

	first = true
	addIndent(buf, indent+2)
	buf.WriteString("node [")
	first = addListAttr(buf, first, "color", s.NodeColor)
	first = addListAttr(buf, first, "style", s.NodeStyle)
	first = addListAttr(buf, first, "shape", s.NodeShape)
	buf.WriteString("];\n")
	first = false

	first = addSubgraphs(buf, first, indent+2, s.Subgraphs)
	first = addNodes(buf, first, indent+2, s.Nodes)
	first = addEdges(buf, first, indent+2, s.Edges)

	addIndent(buf, indent)
	buf.WriteString("}")

	return buf.String()
}
