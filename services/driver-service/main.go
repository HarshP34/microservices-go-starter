package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"ride-sharing/shared/env"
	"ride-sharing/shared/messaging"
	"syscall"

	grpcserver "google.golang.org/grpc"
)

var (
	GrpcAddr = ":9092"
)

func main() {
	rabbitMqURI := env.GetString("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672/")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	svc := NewService()

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

	grpcServer := grpcserver.NewServer()
	NewGRPCHandler(grpcServer, svc)

	log.Printf("Starting grpc server trip service on port: %s", lis.Addr().String())

	// RabbitMQ connection
	rabbitmq, err := messaging.NewRabbitMQ(rabbitMqURI)
	if err != nil {
		log.Fatal(err)
	}

	tripConsumer := NewTripConsumer(rabbitmq, svc)
	go func() {
		if err := tripConsumer.Listen(); err != nil {
			log.Printf("failed to listen for trip events: %v", err)
		}
	}()

	defer rabbitmq.Close()

	log.Println("Starting RabbitMQ connection")

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed to server: %v", err)
			cancel()
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down the server...")
	grpcServer.GracefulStop()
}
