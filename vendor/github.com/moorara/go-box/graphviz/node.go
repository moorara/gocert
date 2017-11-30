package graphviz

import (
	"bytes"
)

// Node represents a graph node for visualization
type Node struct {
	Name      string
	Group     string
	Label     string
	Color     string
	Style     string
	Shape     string
	FontColor string
	FontName  string
}

// NewNode creates a new node for visualization
func NewNode(name, group, label, color, style, shape, fontcolor, fontname string) Node {
	return Node{
		Name:      name,
		Group:     group,
		Label:     label,
		Color:     color,
		Style:     style,
		Shape:     shape,
		FontColor: fontcolor,
		FontName:  fontname,
	}
}

// DotCode generates Graph dot language code for node
func (n *Node) DotCode() string {
	first := true
	buf := new(bytes.Buffer)

	buf.WriteString(n.Name + " [")
	first = addListAttr(buf, first, "group", n.Group)
	first = addListAttr(buf, first, "label", `"`+n.Label+`"`)
	first = addListAttr(buf, first, "color", n.Color)
	first = addListAttr(buf, first, "style", n.Style)
	first = addListAttr(buf, first, "shape", n.Shape)
	first = addListAttr(buf, first, "fontcolor", n.FontColor)
	first = addListAttr(buf, first, "fontname", `"`+n.FontName+`"`)
	buf.WriteString("];")

	return buf.String()
}
