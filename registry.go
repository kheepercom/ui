package ui

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type Registry map[string]Component

func (r Registry) Add(name string, c Component) error {
	name = strings.ToLower(name)
	switch name {
	case "children":
		return errors.New(`"children" is reserved`)
	}
	if _, ok := r[name]; ok {
		return fmt.Errorf("component with name %q already registered", name)
	}
	r[name] = c

	return nil
}

func (reg Registry) Render(r *http.Request, self *html.Node) (*html.Node, error) {
	if self == nil {
		return nil, nil
	}
	nextChild := self.FirstChild
	self.FirstChild = nil
	self.LastChild = nil

	// Render all children first. If this is a custom component the children will
	// be inserted wherever the component has a <children /> element. Otherwise
	// the children will be adopted directly by this node (self).
	var renderedChildren []*html.Node
	for nextChild != nil {
		renderedChild, err := reg.Render(r, nextChild)
		if err != nil {
			return nil, err
		}
		renderedChild.Parent = nil
		renderedChildren = append(renderedChildren, renderedChild)
		nextChild = nextChild.NextSibling
	}

	// If this is a custom component then call its Render method.
	if component, ok := reg[self.Data]; ok && self.Type == html.ElementNode {
		attrs := Attributes{}
		for _, a := range self.Attr {
			if _, ok := attrs[a.Key]; ok {
				attrs[a.Key] = append(attrs[a.Key], a.Val)
			} else {
				attrs[a.Key] = []string{a.Val}
			}
		}
		renderedComponentNode, err := component.Render(r, attrs)
		if err != nil {
			return nil, err
		}
		// The parent of the element with tag "children"
		// TODO error if there are rendered children and nowhere to place them?
		if parent := findParent(renderedComponentNode); parent != nil {
			// This also sets the sibling relations
			placeChildren(parent, renderedChildren)
		}

		return renderedComponentNode, nil
	}

	// Not a component. Adopt the rendered children directly.
	setSiblingRelations(renderedChildren)
	adopt(self, renderedChildren)

	return self, nil
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
