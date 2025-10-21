package main

import (
	"context"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	customerGrpc "test-task/internal/shipment/grpc"
	"test-task/internal/shipment/http"
	"test-task/internal/shipment/repo"
	"test-task/internal/shipment/service"
	"test-task/internal/telemetry"
)

type app struct {
	fiber    *fiber.App
	handlers http.IShipmentHandler
}

func main() {
	ctx := context.Background()
	tp, err := initTracer(ctx)
	if err != nil {
		telemetry.Log.Fatal("failed to init tracer", zap.Error(err))
	}
	defer func() { _ = tp.Shutdown(ctx) }()

	otel.SetTracerProvider(tp)

	if err := telemetry.InitLogger(); err != nil {
		log.Fatal(err)
	}
	telemetry.Log.Info("logger initialized")

	db, err := pgx.Connect(context.Background(), os.Getenv("DB"))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(
		os.Getenv("SERVER_HOST"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	custCl := customerGrpc.NewCustomerClient(conn)
	repo := repo.NewShipmentRepo(db)
	svc := service.NewService(repo, custCl)
	h := http.NewHandler(svc)

	server := &app{
		fiber:    fiber.New(),
		handlers: h,
	}

	fiberApp := fiber.New()
	fiberApp.Use(otelfiber.Middleware())
	server.initRoutes()
	err = server.fiber.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func (r *app) initRoutes() {
	r.fiber.Post("/api/v1/shipments", r.handlers.CreateShipment)
	r.fiber.Get("/api/v1/shipment/:id", r.handlers.GetApplication)
}

func initTracer(ctx context.Context) (*sdktrace.TracerProvider, error) {
	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("shipment-service"),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	return tp, nil
}
