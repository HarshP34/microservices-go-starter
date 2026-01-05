package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	// h "ride-sharing/services/trip-service/internal/infrastructure/http"
	"ride-sharing/services/trip-service/internal/infrastructure/grpc"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"

	// "ride-sharing/shared/env"
	// "time"
	grpcserver "google.golang.org/grpc"
)

var (
	// httpAddr = env.GetString("HTTP_ADDR", ":8083")
	GrpcAddr = ":9093"
)

func main() {
	imemRepo := repository.NewInMemTripRepository()
	svc := service.NewTripService(imemRepo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(){
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	lis, err := net.Listen("tcp", GrpcAddr)
	if err != nil {
		log.Fatalf("Fail to listen %v", err)
	}

	grpcServer := grpcserver.NewServer();

	log.Printf("Starting grpc server trip service on port: %s", lis.Addr().String())
	
	go func(){
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed to server: %v", err)
			cancel()
		}
	}()

	grpc.NewGRPCHandler(grpcServer, svc)

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