package transport

import (
	"context"

	"github.com/go-kit/kit/log"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	kommendpoint "microsrv/kommersant/endpoint"
	"microsrv/kommersant/model"
	"microsrv/pb"
	oldcontext "golang.org/x/net/context"
)

type grpcServer struct {
	create grpctransport.Handler
	result grpctransport.Handler
}

// NewGRPCServer makes a set of endpoints available as a gRPC GreeterServer.
func NewGRPCServer(endpoints kommendpoint.Endpoints, logger log.Logger) pb.KommersantServer {
	options := []grpctransport.ServerOption{grpctransport.ServerErrorLogger(logger)}

	return &grpcServer{
		create: grpctransport.NewServer(
			endpoints.CreateEndpoint,
			decodeGRPCCreateRequest,
			encodeGRPCCreateResponse,
			options...,
		),
		result: grpctransport.NewServer(
			endpoints.ResultEndpoint,
			decodeGRPCCreateRequest,
			encodeGRPCCreateResponse,
			options...,
		),
	}
}

// Create implementation of the method of the GreeterService interface.
func (s *grpcServer) Create(ctx oldcontext.Context, req *pb.KommersantRequest) (*pb.KommersantResponse, error) {
	_, res, err := s.create.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*pb.KommersantResponse), nil
}

// Result implementation of the method of the GreeterService interface.
func (s *grpcServer) Result(ctx oldcontext.Context, req *pb.KommersantRequest) (*pb.KommersantResponse, error) {
	_, res, err := s.result.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*pb.KommersantResponse), nil
}

// decodeGRPCGreetingRequest is a transport/grpc.DecodeRequestFunc that converts
// a gRPC greeting request to a user-domain greeting request.
func decodeGRPCCreateRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.KommersantRequest)
	return model.CreateRequest{AdNum: req.AdNum}, nil
}

// encodeGRPCGreetingResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain greeting response to a gRPC greeting response.
func encodeGRPCCreateResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(model.CreateResponse)
	return &pb.KommersantResponse{Status: res.Status}, nil
}
