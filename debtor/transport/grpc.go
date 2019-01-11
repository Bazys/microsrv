package transport

import (
	"context"

	"api/model"

	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/jinzhu/copier"
	"microsrv/debtor/endpoint"
	"microsrv/pb"
	oldcontext "golang.org/x/net/context"
)

type grpcServer struct {
	createDebtor grpctransport.Handler
	getDebtor    grpctransport.Handler
	getAll       grpctransport.Handler
	save         grpctransport.Handler
	delete       grpctransport.Handler
}

// NewGRPCServer makes a set of endpoints available as a gRPC DebtorServer.
func NewGRPCServer(endpoints debtorendpoint.Endpoints, logger log.Logger) pb.DebtorSvcServer {
	options := []grpctransport.ServerOption{grpctransport.ServerErrorLogger(logger)}

	return &grpcServer{
		createDebtor: grpctransport.NewServer(
			endpoints.CreateDebtorEndpoint,
			decodeGRPCCreateDebtor,
			encodeGRPCDebtorResponse,
			options...,
		),
		getDebtor: grpctransport.NewServer(
			endpoints.GetDebtorEndpoint,
			decodeGRPCGetDebtor,
			encodeGRPCDebtorResponse,
			options...,
		),
		getAll: grpctransport.NewServer(
			endpoints.GetAllDebtorsEndpoint,
			decodeGRPCGetAll,
			encodeGRPCGetAll,
			options...,
		),
		save: grpctransport.NewServer(
			endpoints.SaveDebtorEndpoint,
			decodeGRPCSaveDebtor,
			encodeGRPCDebtorResponse,
			options...,
		),
		delete: grpctransport.NewServer(
			endpoints.DeleteDebtorEndpoint,
			decodeGRPCGetDebtor,
			encodeGRPCDeleteDebtor,
			options...,
		),
	}
}

// Create implementation of the method of the GreeterService interface.
func (s *grpcServer) CreateDebtor(ctx oldcontext.Context, req *pb.Debtor) (*pb.DebtorResponse, error) {
	_, res, err := s.createDebtor.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*pb.DebtorResponse), nil
}

// decodeGRPCGreetingRequest is a transport/grpc.DecodeRequestFunc that converts
// a gRPC greeting request to a user-domain greeting request.
func decodeGRPCCreateDebtor(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.Debtor)
	res := model.Debtor{}
	copier.Copy(&res, req)
	return res, nil
}

// encodeGRPCDebtorResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain greeting response to a gRPC greeting response.
func encodeGRPCDebtorResponse(_ context.Context, response interface{}) (interface{}, error) {
	result := response.(model.DebtorResponse)
	res := pb.DebtorResponse{}
	copier.Copy(&res.Debtor, result.Debtor)
	res.Error = result.Err.Error()
	return &res, nil
}

// GetDebtor implementation of the method of the DebtorServer interface.
func (s *grpcServer) GetDebtor(ctx oldcontext.Context, req *pb.DebtorByID) (*pb.DebtorResponse, error) {
	_, res, err := s.getDebtor.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*pb.DebtorResponse), nil
}

// decodeGRPCGetDebtor is a transport/grpc.DecodeRequestFunc that converts
// a gRPC greeting request to a user-domain greeting request.
func decodeGRPCGetDebtor(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.DebtorByID)
	return req, nil
}

// GetAll implementation of the method of the DebtorServer interface.
func (s *grpcServer) GetAll(ctx oldcontext.Context, req *pb.Pagination) (*pb.DebtorsResponse, error) {
	_, response, err := s.getAll.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	res := response.(pb.DebtorsResponse)
	return &res, nil
}

func decodeGRPCGetAll(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.Pagination)
	return req, nil
}

func encodeGRPCGetAll(_ context.Context, response interface{}) (interface{}, error) {
	result := response.(model.DebtorsResponse)
	res := pb.DebtorsResponse{}
	copier.Copy(&res.Debtors, result.Debtors)
	// res.Error = result.Err.Error()
	return res, nil
}

// Create implementation of the method of the GreeterService interface.
func (s *grpcServer) Save(ctx oldcontext.Context, req *pb.UpadateDebtor) (*pb.DebtorResponse, error) {
	_, res, err := s.save.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*pb.DebtorResponse), nil
}

// decodeGRPCSaveDebtor is a transport/grpc.DecodeRequestFunc that converts
// a gRPC greeting request to a user-domain greeting request.
func decodeGRPCSaveDebtor(_ context.Context, grpReq interface{}) (interface{}, error) {
	req := grpReq.(*pb.UpadateDebtor)
	res := map[string]interface{}{}
	res["debtor"] = req.Update
	res["ID"] = req.ID
	return res, nil
}

// Delete implementation of the method of the DebtorServer interface.
func (s *grpcServer) Delete(ctx oldcontext.Context, req *pb.DebtorByID) (*pb.ErrorResponse, error) {
	_, res, err := s.delete.ServeGRPC(ctx, req)
	result := res.(*pb.ErrorResponse)
	return result, err
}

func encodeGRPCDeleteDebtor(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(model.DebtorResponse)
	return res, nil
}
