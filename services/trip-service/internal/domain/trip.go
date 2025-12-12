package domain

import (
	"context"
	// tripTypes "ride-sharing/services/trip-service/pkg/types"
	tripTypes "ride-sharing/services/trip-service/pkg/types"
	"ride-sharing/shared/types"
	// pbd "ride-sharing/shared/proto/driver"
	"go.mongodb.org/mongo-driver/bson/primitive"
	pb "ride-sharing/shared/proto/trip"
)

type TripModal struct {
	ID       primitive.ObjectID
	UserID   string
	Status   string
	RideFare *RideFareModel
	Driver   *pb.TripDriver
}

type TripRepository interface {
	CreateTrip(ctx context.Context, trips *TripModal) (*TripModal, error)
	SaveRideFare(ctx context.Context, f *RideFareModel) error
	GetRideFareByID(ctx context.Context, id string) (*RideFareModel, error)
}
type TripService interface {
	CreateTrip(ctx context.Context, fare *RideFareModel) (*TripModal, error)

	GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*tripTypes.OsrmApiResponse, error)
	EstimatePackagesPriceWithRoute(rout *tripTypes.OsrmApiResponse) []*RideFareModel
	GenerateTripFares(ctx context.Context,
		fares []*RideFareModel,
		userID string,
		route *tripTypes.OsrmApiResponse,

	) ([]*RideFareModel, error)
	GetAndValidateFare(ctx context.Context, fareID, userID string) (*RideFareModel, error)
}
