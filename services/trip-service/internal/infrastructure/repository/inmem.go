package repository

import (
	"context"
	"fmt"
	"ride-sharing/services/trip-service/internal/domain"
)

type inMemTripRepository struct {
	trips map[string]*domain.TripModel
	rideFares map[string]*domain.RideFareModel
}

func NewInMemTripRepository() *inMemTripRepository {
	return  &inMemTripRepository{
		trips: make(map[string]*domain.TripModel),
		rideFares: make(map[string]*domain.RideFareModel),
	}
}

func (r *inMemTripRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
	r.trips[trip.ID.Hex()] = trip
	return trip, nil
}


func (r *inMemTripRepository) SaveRideFare(ctx context.Context, fare *domain.RideFareModel) error {
	r.rideFares[fare.ID.Hex()] = fare
	return nil
}

func (r *inMemTripRepository) GetRideFareByID(ctx context.Context, fareID string) (*domain.RideFareModel, error) {
	fare, exists := r.rideFares[fareID]
	if !exists {
		return nil, fmt.Errorf("ride fare not found")
	}
	return fare, nil
}

