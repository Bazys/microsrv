package kommersantsvc

import (
	"context"
	"fmt"
	"math/rand"

	"microsrv/kommersant/model"

	"github.com/go-kit/kit/log"
)

// Service describe greetings service.
type Service interface {
	Health() bool
	Create(ctx context.Context, ad model.CreateRequest) (model.CreateResponse, error)
	Result(ctx context.Context, ad model.CreateRequest) (model.CreateResponse, error)
}

// New returns a basic Service with all of the expected middlewares wired in.
func New(logger log.Logger) Service {
	var svc Service
	{
		svc = NewBasicService()
		svc = LoggingMiddleware(logger)(svc)
	}
	return svc
}

// NewBasicService returns a naïve, stateless implementation of Service.
func NewBasicService() Service {
	return basicService{}
}

type basicService struct{}

// Health implementation of the Service.
func (s basicService) Health() bool {
	return true
}

// Create implementation of the Service.
func (s basicService) Create(_ context.Context, ad model.CreateRequest) (model.CreateResponse, error) {
	r := rand.New(rand.NewSource(99))
	return model.CreateResponse{Message: "GO-KIT Hello from Create " + ad.AdNum + " " + fmt.Sprintf("%d", r.Uint32()), Status: 200}, nil
}

// Result implementation of the Service.
func (s basicService) Result(_ context.Context, ad model.CreateRequest) (model.CreateResponse, error) {
	r := rand.New(rand.NewSource(99))
	return model.CreateResponse{Message: "GO-KIT Hello from result " + ad.AdNum + " " + fmt.Sprintf("%d", r.Uint32()), Status: 200}, nil
}
