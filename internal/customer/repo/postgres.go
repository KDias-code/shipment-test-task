package repo

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"test-task/internal/telemetry"
	"time"
)

type ICustomerRepo interface {
	Upsert(ctx context.Context, idn string) (string, error)
	GetCustomer(ctx context.Context, id string) (Customer, error)
}

type repo struct {
	db *pgx.Conn
}

func NewCustomerRepo(db *pgx.Conn) ICustomerRepo {
	return &repo{db}
}

type Customer struct {
	Id        string
	Idn       string
	CreatedAt time.Time
}

func (r repo) Upsert(ctx context.Context, idn string) (string, error) {
	telemetry.TraceLogger(ctx).Info("Upsert", zap.String("idn", idn))

	var id string
	err := r.db.QueryRow(ctx,
		`INSERT INTO customers (idn)
             VALUES ($1)
             ON CONFLICT (idn) DO UPDATE SET idn = EXCLUDED.idn
             RETURNING id`, idn).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r repo) GetCustomer(ctx context.Context, id string) (Customer, error) {
	telemetry.TraceLogger(ctx).Info("GetCustomer", zap.String("id", id))

	var customer Customer

	row := r.db.QueryRow(ctx, `
		SELECT id, idn, created_at FROM customers
		WHERE id = $1
		`, id)

	err := row.Scan(
		&customer.Id,
		&customer.Idn,
		&customer.CreatedAt,
	)
	if err != nil {
		return customer, err
	}

	return customer, nil
}
