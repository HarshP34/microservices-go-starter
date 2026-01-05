package main

import (
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"
)

type previewTripRequest struct {
	UserID      string            `json:"userID"`
	Pickup      types.Coordinate  `json:"pickup"`
	Destination types.Coordinate  `json:"destination"`
}

func (p *previewTripRequest) toProto() *pb.PreviewTripRequest{

	pickupCoord := &pb.Cordinate{
		Latitude: p.Pickup.Latitude,
		Longitude: p.Pickup.Longitude,
	}

	destinationCoord := &pb.Cordinate{
		Latitude: p.Destination.Latitude,
		Longitude: p.Destination.Longitude,
	}
	return  &pb.PreviewTripRequest{
		UserID: p.UserID,
		StartLocation: pickupCoord,
		EndLocation: destinationCoord,
	}
}


type startTripRequest struct {
	RideFareID 		string    `json:"rideFareID"`
	UserID          string    `json:"userID"`
}

func (s *startTripRequest) toProto() *pb.CreateTripRequest {
	return  &pb.CreateTripRequest{
		RideFareID: s.RideFareID,
		UserID: s.UserID,
	}
}