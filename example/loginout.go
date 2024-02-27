package main

import (
	"context"
	_ "embed"
	"net/http"
	"strings"

	"github.com/kheepercom/ui"
	"golang.org/x/net/html"
)

//go:embed login.html
var loginHTML string

//go:embed logout.html
var logoutHTML string

type Loginout struct {
	login  *html.Node
	logout *html.Node
}

func NewLoginout() (*Loginout, error) {
	login, err := html.Parse(strings.NewReader(loginHTML))
	if err != nil {
		return nil, err
	}
	logout, err := html.Parse(strings.NewReader(logoutHTML))
	if err != nil {
		return nil, err
	}
	return &Loginout{
		login:  login,
		logout: logout,
	}, nil
}

type userKey struct{}

func WithUser(ctx context.Context, user string) context.Context {
	return context.WithValue(ctx, userKey{}, user)
}

func UserFromContext(ctx context.Context) string {
	s, _ := ctx.Value(userKey{}).(string)
	return s
}

func (l Loginout) Render(r *http.Request, attrs ui.Attributes) (*html.Node, error) {
	user := UserFromContext(r.Context())
	if user == "" {
		return l.login, nil
	}
	return l.logout, nil
}
