package ui

import "golang.org/x/net/html"

func Clone(n *html.Node) *html.Node {
	children := getChildren(n)
	clonedChildren := make([]*html.Node, len(children))
	for i := range children {
		clonedChildren[i] = Clone(children[i])
	}
	setSiblingRelations(clonedChildren)

	clone := &html.Node{
		Parent:      n.Parent,
		PrevSibling: n.PrevSibling,
		NextSibling: n.NextSibling,
		Type:        n.Type,
		Data:        n.Data,
		DataAtom:    n.DataAtom,
		Attr:        n.Attr,
		// ignore Namespace
	}
	// adopt sets FirstChild and LastChild
	adopt(clone, clonedChildren)

	return clone
}

func getChildren(n *html.Node) []*html.Node {
	var children []*html.Node
	if n == nil {
		return children
	}

	next := n.FirstChild
	for next != nil {
		children = append(children, next)
		next = next.NextSibling
	}

	return children
}
