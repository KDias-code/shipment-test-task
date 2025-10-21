package repo

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"test-task/internal/telemetry"
	"time"
)

type StatusTypes string

const (
	StatusCreated StatusTypes = "CREATED"
)

type Shipment struct {
	ID         string
	Route      string
	Price      decimal.Decimal
	Status     StatusTypes
	CustomerID string
	CreatedAt  time.Time
}

type repo struct {
	db *pgx.Conn
}

type ShipmentRepo interface {
	Create(ctx context.Context, data Shipment) (string, error)
	GetShipmentByID(ctx context.Context, id string) (*Shipment, error)
}

func NewShipmentRepo(db *pgx.Conn) ShipmentRepo {
	return &repo{db}
}

func (r repo) Create(ctx context.Context, data Shipment) (string, error) {
	telemetry.TraceLogger(ctx).Info("Create", zap.String("customer_id", data.CustomerID))

	var id string
	err := r.db.QueryRow(ctx,
		`INSERT INTO shipments (route, price, customer_id) 
             VALUES ($1, $2, $3) RETURNING id`, data.Route, data.Price, data.CustomerID).Scan(&id)
	if err != nil {
		telemetry.TraceLogger(ctx).Error(err.Error(), zap.String("customer_id", data.CustomerID))
		return "", err
	}

	return id, nil
}

func (r repo) GetShipmentByID(ctx context.Context, id string) (*Shipment, error) {
	telemetry.TraceLogger(ctx).Info("GetShipmentByID", zap.String("id", id))

	var s Shipment
	err := r.db.QueryRow(ctx, `
		SELECT id, route, price, status, customer_id, created_at
		FROM shipments
		WHERE id = $1
	`, id).Scan(&s.ID, &s.Route, &s.Price, &s.Status, &s.CustomerID, &s.CreatedAt)
	if err != nil {
		telemetry.TraceLogger(ctx).Error(err.Error(), zap.String("id", id))
		return nil, err
	}

	return &s, nil
}
