package ui

import (
	"bytes"
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

		// rendredComponentNode may have component descendents that require rendering
		// or it may itself be another component
		return reg.Render(r, renderedComponentNode)
	}

	// Not a component. Adopt the rendered children directly.
	setSiblingRelations(renderedChildren)
	adopt(self, renderedChildren)

	return self, nil
}

func (reg Registry) CSS() ([]byte, error) {
	var b bytes.Buffer

	for name, component := range reg {
		if c, ok := component.(Styled); ok {
			css, err := c.CSS()
			if err != nil {
				return nil, fmt.Errorf("failed to get css for component %q: %w", name, err)
			}
			_, err = b.WriteString(css)
			if err != nil {
				return nil, err
			}
		}
	}

	return b.Bytes(), nil
}

func (reg Registry) JS() ([]byte, error) {
	var b bytes.Buffer

	for name, component := range reg {
		if c, ok := component.(Scripted); ok {
			js, err := c.JS()
			if err != nil {
				return nil, fmt.Errorf("failed to get js for component %q: %w", name, err)
			}
			_, err = b.WriteString(js)
			if err != nil {
				return nil, err
			}
		}
	}

	return b.Bytes(), nil
}
