package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	customerpb "test-task/api/proto"
	customerGrpc "test-task/internal/customer/grpc"
	customerRepo "test-task/internal/customer/repo"
	"test-task/internal/customer/service"
	"test-task/internal/telemetry"
)

func main() {
	ctx := context.Background()
	shutdown := telemetry.InitTracer(ctx, "customer-service")
	defer shutdown(ctx)

	db, err := pgx.Connect(context.Background(), os.Getenv("DB"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(ctx)

	repo := customerRepo.NewCustomerRepo(db)
	svc := service.NewService(repo)
	grpcSrv := customerGrpc.NewCustomerServer(svc)

	server := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	customerpb.RegisterCustomerServiceServer(server, grpcSrv)

	ln, _ := net.Listen("tcp", ":9090")
	server.Serve(ln)
}
