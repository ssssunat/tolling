package aggservice

import (
	"context"

	"github.com/ssssunat/tolling/types"
)

type Middleware func(Service) Service

type loggingMiddleware struct {
	next Service
}

func newLoggingMiddleware() Middleware {
	return func(next Service) Service {
		return loggingMiddleware{
			next: next,
		}
	}
}

func (mw loggingMiddleware) Aggregate(_ context.Context, dist types.Distance) error {
	return nil
}

func (mw loggingMiddleware) Calculate(_ context.Context, dist int) (*types.Invoice, error) {
	return nil, nil
}

type intrumentationMiddleware struct {
	next Service
}

func newInstrumentationMiddleware() Middleware {
	return func(next Service) Service {
		return intrumentationMiddleware{
			next: next,
		}
	}
}

func (mw intrumentationMiddleware) Aggregate(_ context.Context, dist types.Distance) error {
	return nil
}

func (mw intrumentationMiddleware) Calculate(_ context.Context, dist int) (*types.Invoice, error) {
	return nil, nil
}
