package graphviz

import (
	"bytes"
)

// Edge represents a graph edge for visualization
type Edge struct {
	From      string
	To        string
	EdgeType  string
	EdgeDir   string
	Label     string
	Color     string
	Style     string
	Arrowhead string
}

// NewEdge creates a new edge for visualization
func NewEdge(from, to, edgeType, edgeDir, label, color, style, arrowhead string) Edge {
	return Edge{
		From:      from,
		To:        to,
		EdgeType:  edgeType,
		EdgeDir:   edgeDir,
		Label:     label,
		Color:     color,
		Style:     style,
		Arrowhead: arrowhead,
	}
}

// DotCode generates Graph dot language code for edge
func (e *Edge) DotCode() string {
	first := true
	buf := new(bytes.Buffer)

	buf.WriteString(e.From + " " + e.EdgeType + " " + e.To + " [")
	first = addListAttr(buf, first, "dirType", e.EdgeDir)
	first = addListAttr(buf, first, "label", `"`+e.Label+`"`)
	first = addListAttr(buf, first, "color", e.Color)
	first = addListAttr(buf, first, "style", e.Style)
	first = addListAttr(buf, first, "arrowhead", e.Arrowhead)
	buf.WriteString("];")

	return buf.String()
}
