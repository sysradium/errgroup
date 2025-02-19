package errgroup

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func supressError() Middleware {
	return func(f GFn) GFn {
		return func() error {
			_ = f()
			return nil
		}
	}
}

func allwaysReturnError(err error) Middleware {
	return func(f GFn) GFn {
		return func() error {
			_ = f()
			return err
		}
	}
}

func TestSingleMiddlwareGetsApplied(t *testing.T) {
	g := New([]Middleware{supressError()})
	g.Go(func() error {
		return errors.New("f")
	})

	require.NoError(t, g.Wait())

}

func TestMiddlwaresAreAppliedInOrderRightToLeft(t *testing.T) {
	err := errors.New("some expected error")
	g := New([]Middleware{allwaysReturnError(err), supressError()})
	g.Go(func() error {
		return errors.New("f")
	})

	require.ErrorIs(t, g.Wait(), err)

}
package errgroup

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroup_TryGo(t *testing.T) {
	tests := []struct {
		name        string
		middlewares []Middleware
		fn          GFn
		want        bool
	}{
		{
			name: "No middlewares, function returns nil",
			middlewares: []Middleware{},
			fn: func() error {
				return nil
			},
			want: true,
		},
		{
			name: "No middlewares, function returns error",
			middlewares: []Middleware{},
			fn: func() error {
				return errors.New("error")
			},
			want: false,
		},
		{
			name: "With middleware, function returns nil",
			middlewares: []Middleware{
				func(next GFn) GFn {
					return func() error {
						return next()
					}
				},
			},
			fn: func() error {
				return nil
			},
			want: true,
		},
		{
			name: "With middleware, function returns error",
			middlewares: []Middleware{
				func(next GFn) GFn {
					return func() error {
						return next()
					}
				},
			},
			fn: func() error {
				return errors.New("error")
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New(tt.middlewares)
			got := g.TryGo(tt.fn)
			assert.Equal(t, tt.want, got)
		})
	}
}
