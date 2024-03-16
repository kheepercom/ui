package ui

import "golang.org/x/net/html"

func adopt(parent *html.Node, children []*html.Node) {
	if parent == nil {
		return
	}
	if len(children) == 0 {
		parent.FirstChild = nil
		parent.LastChild = nil
		return
	}
	parent.FirstChild = children[0]
	parent.LastChild = children[len(children)-1]

	for _, child := range children {
		child.Parent = parent
	}
}

func setSiblingRelations(nodes []*html.Node) {
	for i, node := range nodes {
		node.PrevSibling = nil
		node.NextSibling = nil

		if i == 0 {
			continue
		}

		nodes[i-1].NextSibling = node
		node.PrevSibling = nodes[i-1]
	}
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

// Find the parent of the node named "children"
func findParent(self *html.Node) *html.Node {
	if self.Type == html.ElementNode && self.Data == "children" {
		return self.Parent
	}
	next := self.FirstChild
	for next != nil {
		p := findParent(next)
		if p != nil {
			return p
		}
		next = next.NextSibling
	}

	return nil
}

// Replace the node named "children" with the actual children
func placeChildren(parent *html.Node, orphans []*html.Node) {
	var children []*html.Node

	next := parent.FirstChild
	for next != nil {
		if next.Type == html.ElementNode && next.Data == "children" {
			children = append(children, orphans...)
		} else {
			children = append(children, next)
		}
		next = next.NextSibling
	}

	for i, child := range children {
		child.Parent = parent
		if i == 0 {
			parent.FirstChild = child
			continue
		}
		prevSibling := children[i-1]
		prevSibling.NextSibling = child
		child.PrevSibling = prevSibling
	}
	if len(children) > 0 {
		parent.LastChild = children[len(children)-1]
	}
}
