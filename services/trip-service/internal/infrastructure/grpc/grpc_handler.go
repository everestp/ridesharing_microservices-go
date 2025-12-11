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
	t, err := h.service.GetRoute(ctx,pickupCoord , destinationCoord)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal ,"Failed to get route : %v",err)
		}
// 1 .Estimate the ride  fare prices  based on the route ( eg distance)
estimatedFare := h.service.EstimatePackagesPriceWithRoute(t)
//2.Store the ride date for the create trip 
fare ,err := h.service.GenerateTripFares(ctx, estimatedFare, userID )
if err != nil {
	return nil, status.Errorf(codes.Internal ,"Failed to generate Ride Fare: %v",err)
}

	return  &pb.PreviewTripResponse{
		Route: t.ToProto(),
		RideFares: domain.ToRideFaresProto(fare),
	}, nil















}