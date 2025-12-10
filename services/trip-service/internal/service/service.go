package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	repo domain.TripRepository
}

func NewService(repo domain.TripRepository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModal, error) {
	t := &domain.TripModal{
		ID:       primitive.NewObjectID(),
		UserID:   fare.UserID,
		Status:   "Pending",
		RideFare: fare,
	}
	return s.repo.CreateTrip(ctx, t)
}

func (s *service) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*types.OsrmAPIResponse, error) {
	url := fmt.Sprintf("https://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson", pickup.Longitude, pickup.Latitude, destination.Longitude, destination.Latitude)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to  fetch root form OSRM API: %v", err)
	}
	defer resp.Body.Close()
	
	body ,err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read the response: %v", err)
	}
      var routeResp types.OsrmAPIResponse
	  if err := json.Unmarshal(body, &routeResp); err !=nil{
		return  nil, fmt.Errorf("Failed to read the response: %v", err)
	  }
	  return &routeResp ,nil
}
