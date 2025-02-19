package errgroup

import "log/slog"

func SlogLogger(l *slog.Logger) Middleware {
	return func(g GFn) GFn {
		return func() error {
			err := g()
			if err != nil {
				l.Error("func call failed", "error", err)
			} else {
				l.Debug("func call succeeded")
			}
			return err
		}
	}
}
