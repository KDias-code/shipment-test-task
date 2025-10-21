package grpc

import (
	"context"
	"go.uber.org/zap"
	customerpb "test-task/api/proto"
	"test-task/internal/customer/service"
	"test-task/internal/telemetry"
)

type Server struct {
	customerpb.UnimplementedCustomerServiceServer
	svc service.ICustomerService
}

func NewCustomerServer(s service.ICustomerService) *Server {
	return &Server{
		svc: s,
	}
}

func (s Server) UpsertCustomer(ctx context.Context, req *customerpb.UpsertCustomerRequest) (*customerpb.CustomerResponse, error) {
	id, err := s.svc.UpsertCustomer(ctx, req.Idn)
	if err != nil {
		telemetry.TraceLogger(ctx).Error(err.Error(), zap.String("idn", req.Idn))
		return nil, err
	}

	return &customerpb.CustomerResponse{Id: id}, nil
}

func (s Server) GetCustomer(ctx context.Context, req *customerpb.GetCustomerRequest) (*customerpb.CustomerResponse, error) {
	_, _ = s.svc.GetConsumer(ctx, req.Idn)
	//if err != nil {
	//telemetry.TraceLogger(ctx).Error(err.Error(), zap.String("idn", req.Idn))
	//	return nil, err
	//}

	return &customerpb.CustomerResponse{
		Id:        "mock-id",
		Idn:       "mock-idn",
		CreatedAt: "mock-time",
	}, nil
}
