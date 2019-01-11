package debtorservice

import (
	"context"
	"fmt"
	"time"

	"microsrv/model"

	"github.com/go-kit/kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

// LoggingMiddleware takes a logger as a dependency and returns a ServiceMiddleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{next, logger}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

// Health func
func (mw loggingMiddleware) Health() bool {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Health",
			"healthy", true,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.Health()
}

// Create func
func (mw loggingMiddleware) CreateDebtor(ctx context.Context, d model.Debtor) (model.Debtor, error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "CreateDebtor",
			"debtor.name", d.Name,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.CreateDebtor(ctx, d)
}

// GetDebtor func
func (mw loggingMiddleware) GetDebtor(ctx context.Context, id uint32) (model.Debtor, error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "GetDebtor",
			"Debtor.ID", id,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.GetDebtor(ctx, id)
}

// Save func
func (mw loggingMiddleware) Save(ctx context.Context, d model.Debtor, id uint) (model.Debtor, error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Save",
			"Debtor.ID", id,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.Save(ctx, d, id)
}

// GetAll func
func (mw loggingMiddleware) GetAll(ctx context.Context, p model.Pagination) (model.DebtorsResponse, error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "GetAll",
			"Pagination", fmt.Sprintf("%+v", p),
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.GetAll(ctx, p)
}

// Delete func
func (mw loggingMiddleware) Delete(ctx context.Context, id uint) error {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Delete",
			"Debtor.ID", id,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.Delete(ctx, id)
}
