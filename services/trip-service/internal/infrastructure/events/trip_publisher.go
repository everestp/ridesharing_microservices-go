package events

import (
	"context"
	"ride-sharing/shared/messaging"
	amqp "github.com/rabbitmq/amqp091-go"
	
)



type TripEventPublisher struct {
	rabbitmq *messaging.RabbitMQ

}


func NewTripEventPublisher(rabbitmq *messaging.RabbitMQ) *TripEventPublisher {
	return  &TripEventPublisher{
		rabbitmq: rabbitmq,
	}
}


func (p *TripEventPublisher) PublishWithContext(ctx context.Context) error {
return  p.rabbitmq.PublishMessage(ctx, "hello", "Hello world")