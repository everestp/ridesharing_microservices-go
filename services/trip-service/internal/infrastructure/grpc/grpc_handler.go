package grpc

import (
	"context"
	"log"

	"ride-sharing/services/trip-service/internal/domain"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)



type gRPCHandler struct {
	pb.UnimplementedTripServiceServer
	service domain.TripService
}

func NewGRPCHandler(server *grpc.Server , service domain.TripService) *gRPCHandler{
	handler := &gRPCHandler{
		service: service,
	}
	pb.RegisterTripServiceServer(server, handler)
	return  handler
}
func (h *gRPCHandler) CreateTrip( ctx context.Context,  req *pb.CreateTripRequest)(*pb.CreateTripResponse, error) {
	fareID := req.GetRideFareID()
	userID := req.GetUserID()
// fetch and validte the  fare
   rideFare , err := h.service.GetAndValidateFare(ctx, fareID , userID )
   if err != nil {
	return nil, status.Errorf(codes.Internal ,"Failed to Validate the fare: %v",err)
}

 // 2. Call create trip
 trip , err :=  h.service.CreateTrip(ctx, rideFare)
 if err != nil {
	return nil, status.Errorf(codes.Internal ,"Failed to Create the trip: %v",err)
}

 // 3. WE also need to initiliaze any empty drive to the trip
 // 4. Add a comment at the end of the function to published an event on the aysnc comms module
 return   &pb.CreateTripResponse{
 TripID: trip.ID.Hex(),

 } ,nil
}





func (h *gRPCHandler) PreviewTrip( ctx context.Context,  req *pb.PreviewTripRequest) (*pb.PreviewTripResponse, error) {

	pickup :=req.GetStartLocation()
	destination :=req.GetEndLocation()
   pickupCoord := &types.Coordinate{
	Latitude: pickup.Latitude,
	Longitude: pickup.Longitude,
   }
   destinationCoord := &types.Coordinate{
	Latitude: destination.Latitude,
	Longitude: destination.Longitude,
   }
   userID := req.UserID
	route, err := h.service.GetRoute(ctx,pickupCoord , destinationCoord)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal ,"Failed to get route : %v",err)
		}
// 1 .Estimate the ride  fare prices  based on the route ( eg distance)
estimatedFare := h.service.EstimatePackagesPriceWithRoute(route)
//2.Store the ride date for the create trip 
fare ,err := h.service.GenerateTripFares(ctx, estimatedFare, userID ,route)
if err != nil {
	return nil, status.Errorf(codes.Internal ,"Failed to generate Ride Fare: %v",err)
}

	return  &pb.PreviewTripResponse{
		Route: route.ToProto(),
		RideFares: domain.ToRideFaresProto(fare),
	}, nil















}


