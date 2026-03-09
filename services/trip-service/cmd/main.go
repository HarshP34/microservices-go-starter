package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	// h "ride-sharing/services/trip-service/internal/infrastructure/http"
	"ride-sharing/services/trip-service/internal/infrastructure/events"
	"ride-sharing/services/trip-service/internal/infrastructure/grpc"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/shared/env"
	"ride-sharing/shared/messaging"
	"ride-sharing/shared/tracing"

	// "ride-sharing/shared/env"
	// "time"
	grpcserver "google.golang.org/grpc"
)

var (
	// httpAddr = env.GetString("HTTP_ADDR", ":8083")
	GrpcAddr = ":9093"
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
	rabbitMqURI := env.GetString("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672/")
	imemRepo := repository.NewInMemTripRepository()
	svc := service.NewTripService(imemRepo)

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	lis, err := net.Listen("tcp", GrpcAddr)
	if err != nil {
		log.Fatalf("Fail to listen %v", err)
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

	grpcServer := grpcserver.NewServer(tracing.WithTracingInterceptors()...)

	log.Printf("Starting grpc server trip service on port: %s", lis.Addr().String())

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed to server: %v", err)
			cancel()
		}
	}()

	grpc.NewGRPCHandler(grpcServer, svc, publisher)

	<-ctx.Done()
	log.Println("Shutting down the server...")
	grpcServer.GracefulStop()
	// httpHandler:= &h.HttpHandler{
	// 	Service: svc,
	// }

	// mux.HandleFunc("POST /preview", httpHandler.HandlePreviewTrip)

	// server := &http.Server{
	// 	Addr: httpAddr,
	// 	Handler: mux,
	// }

	// err := server.ListenAndServe(); if err != nil {
	// 	log.Printf("HTTP server error: %v", err)
	// }

	// for {
	// 	time.Sleep(time.Second)
	// }
}
