package errgroup

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Group struct {
	*errgroup.Group
	middlewares []Middleware
}

func (g *Group) Go(fn GFn) {
	res := fn
	for i := len(g.middlewares) - 1; i >= 0; i-- {
		res = g.middlewares[i](res)
	}

	g.Group.Go(res)
}

// make a table test using testify for this function AI!
func (g *Group) TryGo(fn GFn) bool {
	res := fn
	for i := len(g.middlewares) - 1; i >= 0; i-- {
		res = g.middlewares[i](res)
	}

	return g.Group.TryGo(res)
}

func New(middlewares []Middleware) *Group {
	return &Group{
		Group:       &errgroup.Group{},
		middlewares: middlewares,
	}
}

func WithContext(ctx context.Context) (Group, context.Context) {
	g, ctx := errgroup.WithContext(ctx)
	return Group{
		Group: g,
	}, ctx
}

type GFn func() error

type Middleware func(GFn) GFn
