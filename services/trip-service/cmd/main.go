package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/danielmoisa/bolt-app/pkg/db"
	"github.com/danielmoisa/bolt-app/pkg/env"
	"github.com/danielmoisa/bolt-app/pkg/messaging"
	"github.com/danielmoisa/bolt-app/pkg/tracing"

	"github.com/danielmoisa/bolt-app/services/trip-service/internal/infrastructure/events"
	"github.com/danielmoisa/bolt-app/services/trip-service/internal/infrastructure/grpc"
	"github.com/danielmoisa/bolt-app/services/trip-service/internal/infrastructure/repository"
	"github.com/danielmoisa/bolt-app/services/trip-service/internal/service"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	grpcserver "google.golang.org/grpc"
)

var (
	GrpcAddr    = ":9093"
	MetricsAddr = ":9193"
)

func main() {
	// Initialize Tracing
	tracerCfg := tracing.Config{
		ServiceName:    "trip-service",
		Environment:    env.GetString("ENVIRONMENT", "development"),
		JaegerEndpoint: env.GetString("JAEGER_ENDPOINT", "http://jaeger:14268/api/traces"),
	}

	sh, err := tracing.InitTracer(tracerCfg)
	if err != nil {
		log.Fatalf("Failed to initialize the tracer: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer sh(ctx)

	// Initialize Postgres DB
	pgClient, err := db.NewPostgresClient(ctx, db.NewPostgresDefaultConfig())
	if err != nil {
		log.Fatalf("Failed to initialize PostgresDB, err: %v", err)
	}
	defer pgClient.Close()

	pg := db.GetDB(pgClient)

	rabbitMqURI := env.GetString("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672/")

	pgRepo := repository.NewTripRepository(pg)
	svc := service.NewService(pgRepo)

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	lis, err := net.Listen("tcp", GrpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// RabbitMQ connection
	rabbitmq, err := messaging.NewRabbitMQ(rabbitMqURI)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitmq.Close()

	log.Println("Starting RabbitMQ connection")

	publisher := events.NewTripEventPublisher(rabbitmq)

	// Start driver consumer
	driverConsumer := events.NewDriverConsumer(rabbitmq, svc)
	go driverConsumer.Listen()

	// Start payment consumer
	paymentConsumer := events.NewPaymentConsumer(rabbitmq, svc)
	go paymentConsumer.Listen()

	// Starting the gRPC server
	grpcServer := grpcserver.NewServer(tracing.WithTracingInterceptors()...)
	grpc.NewGRPCHandler(grpcServer, svc, publisher)

	log.Printf("Starting gRPC server Trip service on port %s", lis.Addr().String())

	// Start metrics server
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
		log.Printf("Starting metrics server on %s", MetricsAddr)
		if err := http.ListenAndServe(MetricsAddr, mux); err != nil {
			log.Printf("Failed to start metrics server: %v", err)
		}
	}()

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
			cancel()
		}
	}()

	// wait for the shutdown signal
	<-ctx.Done()
	log.Println("Shutting down the server...")
	grpcServer.GracefulStop()
}
