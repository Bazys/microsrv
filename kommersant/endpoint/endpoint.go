package kommendpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kommersantmodel "microsrv/kommersant/model"
	kommersantsvc "microsrv/kommersant/service"
)

// Endpoints collects all of the endpoints that compose a greeter service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	HealthEndpoint endpoint.Endpoint // used by Consul for the healthcheck
	CreateEndpoint endpoint.Endpoint
	ResultEndpoint endpoint.Endpoint
}

// MakeServerEndpoints returns service Endoints, and wires in all the provided
// middlewares.
func MakeServerEndpoints(s kommersantsvc.Service, logger log.Logger) Endpoints {
	var healthEndpoint endpoint.Endpoint
	{
		healthEndpoint = MakeHealthEndpoint(s)
		healthEndpoint = LoggingMiddleware(log.With(logger, "method", "Health"))(healthEndpoint)
	}

	var createEndpoint endpoint.Endpoint
	{
		createEndpoint = MakeCreateEndpoint(s)
		createEndpoint = LoggingMiddleware(log.With(logger, "method", "Create"))(createEndpoint)
	}
	var resultEndpoint endpoint.Endpoint
	{
		resultEndpoint = MakeResultEndpoint(s)
		resultEndpoint = LoggingMiddleware(log.With(logger, "method", "Result"))(resultEndpoint)
	}

	return Endpoints{
		HealthEndpoint: healthEndpoint,
		CreateEndpoint: createEndpoint,
		ResultEndpoint: resultEndpoint,
	}
}

// MakeHealthEndpoint constructs a Health endpoint wrapping the service.
func MakeHealthEndpoint(s kommersantsvc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		healthy := s.Health()
		return kommersantmodel.HealthResponse{Healthy: healthy}, nil
	}
}

// MakeCreateEndpoint constructs a Greeter endpoint wrapping the service.
func MakeCreateEndpoint(s kommersantsvc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(kommersantmodel.CreateRequest)
		return s.Create(ctx, req)
	}
}

// MakeResultEndpoint constructs a Greeter endpoint wrapping the service.
func MakeResultEndpoint(s kommersantsvc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(kommersantmodel.CreateRequest)
		return s.Result(ctx, req)
	}
}

// Failer is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so if they've
// failed, and if so encode them using a separate write path based on the error.
type Failer interface {
	Failed() error
}
