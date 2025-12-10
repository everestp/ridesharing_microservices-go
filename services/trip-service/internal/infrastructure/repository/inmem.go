package repository

import (
	"context"
	"ride-sharing/services/trip-service/internal/domain"
)

type inmemRepository struct {
	trips map[string] *domain.TripModal
	rideFare map[string] *domain.RideFareModel
}

func NewInmemRepository() *inmemRepository {
	return  &inmemRepository{
		trips: make(map[string]*domain.TripModal),
		rideFare: make(map[string]*domain.RideFareModel),
	}
}



func (r *inmemRepository)CreateTrip( ctx context.Context , trip *domain.TripModal) (*domain.TripModal , error){

	r.trips[trip.ID.Hex()]= trip
	return trip, nil

}