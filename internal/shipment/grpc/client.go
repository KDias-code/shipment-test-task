package grpc

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	customerpb "test-task/api/proto"
	"test-task/internal/telemetry"
)

type ICustomerClient interface {
	UpsertCustomer(ctx context.Context, idn string) (string, error)
	GetCustomer(ctx context.Context, idn string) (*CustomerResponse, error)
}

type CustomerClient struct {
	client customerpb.CustomerServiceClient
}

func NewCustomerClient(conn *grpc.ClientConn) ICustomerClient {
	return &CustomerClient{
		client: customerpb.NewCustomerServiceClient(conn),
	}
}

func (c *CustomerClient) UpsertCustomer(ctx context.Context, idn string) (string, error) {
	resp, err := c.client.UpsertCustomer(ctx, &customerpb.UpsertCustomerRequest{Idn: idn})
	if err != nil {
		telemetry.TraceLogger(ctx).Error(err.Error(), zap.String("idn", idn))
		return "", err
	}

	return resp.Id, nil
}

func (c *CustomerClient) GetCustomer(ctx context.Context, idn string) (*CustomerResponse, error) {
	resp, err := c.client.GetCustomer(ctx, &customerpb.GetCustomerRequest{Idn: idn})
	if err != nil {
		telemetry.TraceLogger(ctx).Error(err.Error(), zap.String("idn", idn))
		return nil, err
	}

	return &CustomerResponse{
		ID:        resp.Id,
		IDN:       resp.Idn,
		CreatedAt: resp.CreatedAt,
	}, nil
}

type CustomerResponse struct {
	ID        string
	IDN       string
	CreatedAt string
}
