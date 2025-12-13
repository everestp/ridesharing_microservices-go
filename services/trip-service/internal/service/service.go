package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	tripTypes "ride-sharing/services/trip-service/pkg/types"
		// tripTypes "ride-sharing/services/trip-service/pkg/types"
	// pbd "ride-sharing/shared/proto/driver"
	pb "ride-sharing/shared/proto/trip"
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
		Driver: &pb.TripDriver{},
		
	}
	return s.repo.CreateTrip(ctx, t)
}

func (s *service) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*tripTypes.OsrmApiResponse, error) {
	url := fmt.Sprintf("https://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson", pickup.Longitude, pickup.Latitude, destination.Longitude, destination.Latitude)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to  fetch root form OSRM API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read the response: %v", err)
	}
	var routeResp tripTypes.OsrmApiResponse
	if err := json.Unmarshal(body, &routeResp); err != nil {
		return nil, fmt.Errorf("Failed to read the response: %v", err)
	}
	return &routeResp, nil
}

func (s *service) EstimatePackagesPriceWithRoute(rout *tripTypes.OsrmApiResponse) []*domain.RideFareModel {
	baseFares := getBaseFares()
	estimatedFare := make([]*domain.RideFareModel, len(baseFares))
	for i, f := range baseFares {
		estimatedFare[i] = estimatedFareRoute(f, rout)
	}
	return estimatedFare
}




func (s *service) GetAndValidateFare(ctx context.Context , fareID , userID  string, ) (*domain.RideFareModel , error){
  fare , err := s.repo.GetRideFareByID(ctx, fareID)
 if err != nil{
		return  nil , fmt.Errorf("Failed to get Trip Fare : %w",err)
	}
	if  fare == nil {
		return nil , fmt.Errorf("Fare  does not exist")
	}

	// User fare validation (user is owner  of this  fare?)
	if userID  !=fare.UserID {
		return nil , fmt.Errorf("Fare not belong to the  user : v")
	}
	return fare, nil
}
func  estimatedFareRoute(f *domain.RideFareModel, rout *tripTypes.OsrmApiResponse) *domain.RideFareModel {
	pricingCfg := tripTypes.DefaultPricingConfig()
	carPackagePrice := f.TotalPriceInCents
	distanceKm := rout.Routes[0].Distance
	durationInMinutes := rout.Routes[0].Duration

	// distance 
	distaceFare := distanceKm *pricingCfg.PricePerUnitOfDistance

	//time
	timeFare := durationInMinutes * pricingCfg.PricingPerMinute
	// car price
 totaPrice := carPackagePrice + distaceFare + timeFare
 return  &domain.RideFareModel{
	TotalPriceInCents: totaPrice,
	PackageSlug: f.PackageSlug,
 }

}
func (s *service) GenerateTripFares(ctx context.Context, ridefare []*domain.RideFareModel, userID string , route *tripTypes.OsrmApiResponse) ([]*domain.RideFareModel, error) {  
  fares := make([]*domain.RideFareModel , len(ridefare))
  for i ,f := range ridefare {
	id := primitive.NewObjectID()
	fare := &domain.RideFareModel{
		UserID: userID,
		ID: id,
		TotalPriceInCents: f.TotalPriceInCents,
		PackageSlug: f.PackageSlug,
		Route: route,

	}

	if err := s.repo.SaveRideFare(ctx ,fare); err != nil{
		return  nil , fmt.Errorf("Failed to saved Trip  Fare : %w",err)
	}
	fares[i]=fare;
}
return  fares ,nil

}

func getBaseFares() []*domain.RideFareModel {
	return []*domain.RideFareModel{
		{
			PackageSlug: "suv",
			TotalPriceInCents:  200,
		},
		{
			PackageSlug: "sedan",
			TotalPriceInCents:  500,
		},
		{
			PackageSlug: "luxury",
			TotalPriceInCents:  1000,
		},
	}
}

func (s *service) GetTripByID(ctx context.Context, id string) (*domain.TripModel, error) {
	return s.repo.GetTripByID(ctx, id)
}

func (s *service) UpdateTrip(ctx context.Context, tripID string, status string, driver *pbd.Driver) error {
	return s.repo.UpdateTrip(ctx, tripID, status, driver)
}
