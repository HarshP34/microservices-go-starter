package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"

	tripTypes "ride-sharing/services/trip-service/pkg/types"
)

type service struct {
	repo domain.TripRepository
}

func NewTripService(repo domain.TripRepository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error) {
	trip := &domain.TripModel{
		ID:      primitive.NewObjectID(),
		UserID:  fare.UserID,
		Status:  "created",
		RideFare: fare,
		Driver: &trip.TripDriver{},
	}
	return s.repo.CreateTrip(ctx, trip)
}


func (s *service) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*tripTypes.OsrmApiResponse, error) {
	url := fmt.Sprintf(
		"http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson",
	pickup.Longitude, pickup.Latitude, destination.Longitude, destination.Latitude,
)

	resp, err := http.Get(url)
		if err != nil {
		return nil, fmt.Errorf("failed to fetch route from OSRM API: %v", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var routeResp tripTypes.OsrmApiResponse
	err = json.Unmarshal(respBody, &routeResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the response: %v", err)
	}
	return &routeResp, nil
}

func (s *service) EstimatePackagesPriceWithRoute(route *tripTypes.OsrmApiResponse) ([]*domain.RideFareModel, error) {
	baseFares := getBaseFares()
	estimatedFares := make([]*domain.RideFareModel, len(baseFares))

	for i, f := range baseFares {
		estimatedFares[i] = estimateFareRoute(f, route)
	}

	return estimatedFares, nil
}


func (s *service) GenerateTripFares(ctx context.Context, rideFares []*domain.RideFareModel, userID string, route *tripTypes.OsrmApiResponse) ([]*domain.RideFareModel, error) {
	fares := make([]*domain.RideFareModel, len(rideFares))
	for i, f := range rideFares {
		id := primitive.NewObjectID()
		fare := &domain.RideFareModel{
			UserID:       	   userID,
			ID:           	   id,
			TotalPriceInCents: f.TotalPriceInCents,
			PackageSlug:       f.PackageSlug,
			Route:             route,
		}

		if err := s.repo.SaveRideFare(ctx, fare); err != nil {
			return nil, fmt.Errorf("failed to save trip fare: %v", err)
		}

		fares[i] = fare
	}

	return fares, nil
}

func (s *service) GetAndValidateFare(ctx context.Context, fareID, userID string) (*domain.RideFareModel, error) {
	fare, err := s.repo.GetRideFareByID(ctx, fareID)
	if err != nil {
		return  nil, fmt.Errorf("failed to get trip fare: %w", err)
	}

	if fare != nil {
		return  nil, fmt.Errorf("trip fare is not exits: %w", err)
	}

	if userID != fare.UserID {
		return  nil, fmt.Errorf("trip fare does not belong to user: %s", userID)
	}

	return fare, nil
}



func estimateFareRoute(f *domain.RideFareModel, route *tripTypes.OsrmApiResponse) *domain.RideFareModel {
	pricingCfg := tripTypes.DefaultPricingConfig()
	carPackagePrice := f.TotalPriceInCents

	distanceKm := route.Routes[0].Distance
	durateInMinutes := route.Routes[0].Distance

	distanceFare := distanceKm * pricingCfg.PricePerUnitOfDistance
	timeFare := durateInMinutes * pricingCfg.PricingPerMinutes
	totalPrice := carPackagePrice + distanceFare + timeFare

	return &domain.RideFareModel {
		TotalPriceInCents: totalPrice,
		PackageSlug: f.PackageSlug,
	}
}

func getBaseFares() []*domain.RideFareModel {
	return []*domain.RideFareModel{
		{
			PackageSlug:       "sedan",
			TotalPriceInCents: 200,
		},
		{
			PackageSlug:       "suv",
			TotalPriceInCents: 350,
		},
		{
			PackageSlug:       "van",
			TotalPriceInCents: 400,
		},
		{
			PackageSlug:       "luxury",
			TotalPriceInCents: 1000,
		},
	}
}