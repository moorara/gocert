package graphviz

import (
	"bytes"
)

func addIndent(buf *bytes.Buffer, indent int) {
	for i := 0; i < indent; i++ {
		buf.WriteString(" ")
	}
}

func addListAttr(buf *bytes.Buffer, first bool, name, value string) bool {
	if value != "" && value != "\"\"" {
		if !first {
			buf.WriteString(", ")
		}
		buf.WriteString(name)
		buf.WriteString("=")
		buf.WriteString(value)
		return false
	}

	return first
}

func addAttr(buf *bytes.Buffer, first bool, indent int, name, value string) bool {
	if value != "" && value != "\"\"" {
		addIndent(buf, indent)
		buf.WriteString(name)
		buf.WriteString("=")
		buf.WriteString(value)
		buf.WriteString(";\n")
		return false
	}

	return first
}

func addNodes(buf *bytes.Buffer, first bool, indent int, nodes []Node) bool {
	added := len(nodes) > 0
	if !first && added {
		buf.WriteString("\n")
	}

	for _, node := range nodes {
		addIndent(buf, indent)
		buf.WriteString(node.DotCode())
		buf.WriteString("\n")
	}

	return first && !added
}

func addEdges(buf *bytes.Buffer, first bool, indent int, edges []Edge) bool {
	added := len(edges) > 0
	if !first && added {
		buf.WriteString("\n")
	}

	for _, edge := range edges {
		addIndent(buf, indent)
		buf.WriteString(edge.DotCode())
		buf.WriteString("\n")
	}

	return first && !added
}

func addSubgraphs(buf *bytes.Buffer, first bool, indent int, subgraphs []Subgraph) bool {
	added := len(subgraphs) > 0

	for _, subgraph := range subgraphs {
		if !first {
			buf.WriteString("\n")
		}
		buf.WriteString(subgraph.DotCode(indent))
		buf.WriteString("\n")
		first = false
	}

	return first && !added
}
