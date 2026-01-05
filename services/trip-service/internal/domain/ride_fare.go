package domain

import (
	pb "ride-sharing/shared/proto/trip"

	tripTypes "ride-sharing/services/trip-service/pkg/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RideFareModel struct {
	ID        		  primitive.ObjectID
	UserID    		  string
	PackageSlug 	  string
	TotalPriceInCents float64
	Route             *tripTypes.OsrmApiResponse
}

func (r *RideFareModel) ToProto() *pb.RideFare {
	return  &pb.RideFare{
		Id: r.ID.Hex(),
		UserID: r.UserID,
		TotalPriceInCents: r.TotalPriceInCents,
		PackageSlug: r.PackageSlug,
	}
}

func ToRideFaresProto(fares []*RideFareModel) []*pb.RideFare {
	pbFares := make([]*pb.RideFare, len(fares))

	for i, fare := range fares {
		pbFares[i] = fare.ToProto()
	}

	return pbFares
}