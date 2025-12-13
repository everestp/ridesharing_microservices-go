package events

import (
	"context"
	"encoding/json"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"
	// amqp "github.com/rabbitmq/amqp091-go"
)



type TripEventPublisher struct {
	rabbitmq *messaging.RabbitMQ

}


func NewTripEventPublisher(rabbitmq *messaging.RabbitMQ) *TripEventPublisher {
	return  &TripEventPublisher{
		rabbitmq: rabbitmq,
	}
}
func (p *TripEventPublisher) PublishTripCreated(ctx context.Context, trip *domain.TripModal) error {
	payload := messaging.TripEventData{
		Trip: trip.ToProto(),
	}

	tripEventJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return p.rabbitmq.PublishMessage(ctx, contracts.TripEventCreated, contracts.AmqpMessage{
		OwnerID: trip.UserID,
		Data:    tripEventJSON,
	})
}



func (p *TripEventPublisher) PublishWithContext(ctx context.Context , trip *domain.TripModal) error {

	tripEventJSON , err := json.Marshal(trip)
	if err != nil {
		return  err
	}
return  p.rabbitmq.PublishMessage(ctx,contracts.TripEventCreated, contracts.AmqpMessage{
	OwnerID: trip.UserID,
	Data: tripEventJSON,

})
}