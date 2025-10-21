package service

import (
	"context"
	"go.uber.org/zap"
	"test-task/internal/shipment/grpc"
	"test-task/internal/shipment/repo"
	"test-task/internal/telemetry"
)

type CreateShipmentResponse struct {
	CustId     string
	Status     repo.StatusTypes
	ShipmentId string
}

type IShipmentService interface {
	Create(ctx context.Context, shipment repo.Shipment, idn string) (*CreateShipmentResponse, error)
	GetShipment(ctx context.Context, id string) (*repo.Shipment, error)
}

type Service struct {
	repo    repo.ShipmentRepo
	grpcCli grpc.ICustomerClient
}

func NewService(repo repo.ShipmentRepo, grpcCli grpc.ICustomerClient) IShipmentService {
	return &Service{repo, grpcCli}
}

func (s Service) Create(ctx context.Context, shipment repo.Shipment, idn string) (*CreateShipmentResponse, error) {
	telemetry.TraceLogger(ctx).Info("creating shipment", zap.String("idn", idn))

	custID, err := s.grpcCli.UpsertCustomer(ctx, idn)
	if err != nil {
		return nil, err
	}

	shipment.CustomerID = custID
	shipment.Status = repo.StatusCreated
	shipId, err := s.repo.Create(ctx, shipment)
	if err != nil {
		return nil, err
	}

	return &CreateShipmentResponse{
		CustId:     custID,
		Status:     shipment.Status,
		ShipmentId: shipId,
	}, nil
}

func (s Service) GetShipment(ctx context.Context, id string) (*repo.Shipment, error) {
	telemetry.TraceLogger(ctx).Info("get shipment", zap.String("id", id))

	_, err := s.grpcCli.GetCustomer(ctx, id)
	if err != nil {
		telemetry.TraceLogger(ctx).Error(err.Error(), zap.String("id", id))
		return nil, err
	}

	shipment, err := s.repo.GetShipmentByID(ctx, id)
	if err != nil {
		telemetry.TraceLogger(ctx).Error(err.Error(), zap.String("id", id))
		return nil, err
	}

	return shipment, nil
}
