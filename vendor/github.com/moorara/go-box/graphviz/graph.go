package graphviz

import (
	"bytes"
)

// Graph represents a graph for visualization
type Graph struct {
	Strict    bool
	Digraph   bool
	Name      string
	Rankdir   string
	NodeColor string
	NodeStyle string
	NodeShape string
	Nodes     []Node
	Edges     []Edge
	Subgraphs []Subgraph
}

// NewGraph creates a new graph for visualization
func NewGraph(strict, digraph bool, name, rankdir, nodeColor, nodeStyle, nodeShape string) Graph {
	return Graph{
		Strict:    strict,
		Digraph:   digraph,
		Name:      name,
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
func (g *Graph) AddNode(nodes ...Node) {
	for _, n := range nodes {
		g.Nodes = append(g.Nodes, n)
	}
}

// AddEdge adds a new edge to graph for visualization
func (g *Graph) AddEdge(edges ...Edge) {
	for _, e := range edges {
		g.Edges = append(g.Edges, e)
	}
}

// AddSubgraph adds a new subgraph to graph for visualization
func (g *Graph) AddSubgraph(subgraphs ...Subgraph) {
	for _, sg := range subgraphs {
		g.Subgraphs = append(g.Subgraphs, sg)
	}
}

// DotCode generates Graph dot language code for graph
func (g *Graph) DotCode() string {
	first := true
	buf := new(bytes.Buffer)

	if g.Strict {
		buf.WriteString("strict ")
	}

	if g.Digraph {
		buf.WriteString("digraph ")
	} else {
		buf.WriteString("graph ")
	}

	if g.Name != "" {
		buf.WriteString(g.Name)
		buf.WriteString(" ")
	}

	buf.WriteString("{\n")

	first = addAttr(buf, first, 2, "rankdir", g.Rankdir)

	first = true
	addIndent(buf, 2)
	buf.WriteString("node [")
	first = addListAttr(buf, first, "color", g.NodeColor)
	first = addListAttr(buf, first, "style", g.NodeStyle)
	first = addListAttr(buf, first, "shape", g.NodeShape)
	buf.WriteString("];\n")
	first = false

	first = addSubgraphs(buf, first, 2, g.Subgraphs)
	first = addNodes(buf, first, 2, g.Nodes)
	addEdges(buf, first, 2, g.Edges)
	buf.WriteString("}")

	return buf.String()
}
