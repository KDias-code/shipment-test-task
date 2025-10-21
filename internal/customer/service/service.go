package service

import (
	"context"
	"test-task/internal/customer/repo"
)

type ICustomerService interface {
	UpsertCustomer(ctx context.Context, idn string) (string, error)
	GetConsumer(ctx context.Context, idn string) (*repo.Customer, error)
}

type CustomerService struct {
	repo repo.ICustomerRepo
}

func NewService(repo repo.ICustomerRepo) *CustomerService { return &CustomerService{repo: repo} }

func (s *CustomerService) UpsertCustomer(ctx context.Context, idn string) (string, error) {
	res, err := s.repo.Upsert(ctx, idn)
	if err != nil {
		return "", err
	}

	return res, nil
}

func (s *CustomerService) GetConsumer(ctx context.Context, idn string) (*repo.Customer, error) {
	res, err := s.repo.GetCustomer(ctx, idn)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
