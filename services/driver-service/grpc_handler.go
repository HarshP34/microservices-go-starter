package main

import (
	"context"
	pb "ride-sharing/shared/proto/driver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCHandler struct {
	pb.UnimplementedDriverServiceServer
	service *Service
}

func NewGRPCHandler(server *grpc.Server, service *Service) *gRPCHandler {
	handler := &gRPCHandler{
		service: service,
	}
	pb.RegisterDriverServiceServer(server, handler)
	return handler
}

func (g *gRPCHandler) RegisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	driver, err := g.service.RegisterDriver(req.GetDriverID(), req.GetPackageSlug())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register driver: %v", err)
	}
	return &pb.RegisterDriverResponse{
		Driver: driver,
	}, nil
}

func (g *gRPCHandler) UnRegisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	err := g.service.UnRegisterDriver(req.GetDriverID())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unregister driver: %v", err)
	}
	return &pb.RegisterDriverResponse{
		Driver: &pb.Driver{
			Id: req.GetDriverID(),
		},
	}, nil
}
