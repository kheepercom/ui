package ui

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func TestParse(t *testing.T) {
	type test struct {
		in      string
		want    *html.Node
		wantErr bool
	}
	tests := map[string]test{
		"<div></div>": {
			in:   "<div></div>",
			want: &html.Node{Type: html.ElementNode, Data: "div", DataAtom: atom.Div},
		},
		"<div />": {
			in:   "<div />",
			want: &html.Node{Type: html.ElementNode, Data: "div", DataAtom: atom.Div},
		},
		"<MyDiv></MyDiv>": {
			in:   "<MyDiv></MyDiv>",
			want: &html.Node{Type: html.ElementNode, Data: "mydiv"},
		},
		"<MyDiv />": {
			in:   "<MyDiv />",
			want: &html.Node{Type: html.ElementNode, Data: "mydiv"},
		},
		"<oops": {
			in:      "<oops",
			wantErr: true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := Parse(strings.NewReader(test.in))
			if test.wantErr {
				require.Error(t, err)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.want, got)
			}
		})
	}
}
