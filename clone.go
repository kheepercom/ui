package ui

import "golang.org/x/net/html"

func Clone(n *html.Node) *html.Node {
	children := getChildren(n)
	clonedChildren := make([]*html.Node, len(children))
	for i := range children {
		clonedChildren[i] = Clone(children[i])
	}
	setSiblingRelations(clonedChildren)

	// FirstChild and LastChild are set in adopt
	clone := &html.Node{
		Parent:      n.Parent,
		PrevSibling: n.PrevSibling,
		NextSibling: n.NextSibling,
		Type:        n.Type,
		Data:        n.Data,
		DataAtom:    n.DataAtom,
		Attr:        n.Attr,
		Namespace:   n.Namespace,
	}
	adopt(clone, clonedChildren)

	return clone
}
