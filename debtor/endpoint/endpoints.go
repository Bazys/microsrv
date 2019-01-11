package debtorendpoint

import (
	"context"

	"microsrv/pb"

	"microsrv/model"

	"microsrv/debtor/service"

	"github.com/go-kit/kit/endpoint"
	"github.com/jinzhu/copier"
)

// Endpoints struct
type Endpoints struct {
	HealthEndpoint        endpoint.Endpoint // used by Consul for the healthcheck
	CreateDebtorEndpoint  endpoint.Endpoint
	GetDebtorEndpoint     endpoint.Endpoint
	GetAllDebtorsEndpoint endpoint.Endpoint
	SaveDebtorEndpoint    endpoint.Endpoint
	DeleteDebtorEndpoint  endpoint.Endpoint
}

// MakeServerEndpoints func
func MakeServerEndpoints(s debtorservice.Service) Endpoints {
	return Endpoints{
		HealthEndpoint:        HealthEndpoint(s),
		CreateDebtorEndpoint:  CreateEndpoint(s),
		GetDebtorEndpoint:     GetEndpoint(s),
		GetAllDebtorsEndpoint: GetAllEndpoint(s),
		SaveDebtorEndpoint:    SaveEndpoint(s),
		DeleteDebtorEndpoint:  DeleteEndpoint(s),
	}
}

// compile time assertions for our response types implementing endpoint.Failer.
var (
	_ endpoint.Failer = HealthResponse{}
	_ endpoint.Failer = model.DebtorsResponse{}
	_ endpoint.Failer = model.DebtorResponse{}
)

// HealthEndpoint constructs a Health endpoint wrapping the service.
func HealthEndpoint(s debtorservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		healthy := s.Health()
		return HealthResponse{Healthy: healthy}, nil
	}
}

// CreateEndpoint func
func CreateEndpoint(s debtorservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(model.Debtor)
		return s.CreateDebtor(ctx, req)
	}
}

// SaveEndpoint func
func SaveEndpoint(s debtorservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := response.(map[string]interface{})
		return s.Save(ctx, req["debtor"].(model.Debtor), req["ID"].(uint))
	}
}

// GetEndpoint func
func GetEndpoint(s debtorservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.DebtorByID)
		res, e := s.GetDebtor(ctx, req.ID)
		resp := model.DebtorResponse{}
		resp.Debtor = res
		resp.Err = e
		return resp, nil
	}
}

// GetAllEndpoint func
func GetAllEndpoint(s debtorservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.Pagination)
		page := model.Pagination{}
		copier.Copy(&page, req)
		resp := model.DebtorsResponse{}
		resp, e := s.GetAll(ctx, page)
		resp.Err = e
		return resp, nil
	}
}

// DeleteEndpoint func
func DeleteEndpoint(s debtorservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.DebtorByID)
		return nil, s.Delete(ctx, uint(req.ID))
	}
}

// Failer is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so if they've
// failed, and if so encode them using a separate write path based on the error.
type Failer interface {
	Failed() error
}

// HealthRequest collects the request parameters for the Health method.
type HealthRequest struct{}

// HealthResponse collects the response values for the Health method.
type HealthResponse struct {
	Healthy bool  `json:"healthy,omitempty"`
	Err     error `json:"err,omitempty"`
}

// Failed implements Failer.
func (r HealthResponse) Failed() error { return r.Err }
