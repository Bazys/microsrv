package kommersantsvc

import (
	"context"
	"time"

	"microsrv/kommersant/model"

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
func (mw loggingMiddleware) Create(ctx context.Context, ad model.CreateRequest) (model.CreateResponse, error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Create",
			"ad_num", ad.AdNum,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.Create(ctx, ad)
}

// Result func
func (mw loggingMiddleware) Result(ctx context.Context, ad model.CreateRequest) (model.CreateResponse, error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Result",
			"ad_num", ad.AdNum,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.Result(ctx, ad)
}
