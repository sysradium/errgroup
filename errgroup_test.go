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
