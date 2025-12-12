package messaging

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	Channel *amqp.Channel
}



func NewRabbitMQ(uri string) (*RabbitMQ, error) {
	//RabbitMQ connection
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("Failed to Connect to Rabbit mq : %v", err)

	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("Failed to create Channel : %v", err)
	}
	rmq := &RabbitMQ{
		conn:    conn,
		Channel: ch,
	}
	if err := rmq.setupExchangesAndQueue(); err != nil {
		//Clean up iof setup fails
          rmq.Close()
		 return nil, fmt.Errorf("Failed to Setup exchanges and Queue:  %v", err)
	}
	return rmq, nil
}

func (r *RabbitMQ) PublishMessage(ctx context.Context , routingKey string , message string) error{
 return r.Channel.PublishWithContext(ctx,
  "",     // exchange
 "hello", // routing key
  false,  // mandatory
  false,  // immediate
  amqp.Publishing {
    ContentType: "text/plain",
    Body:        []byte(message),
	DeliveryMode: amqp.Persistent ,
  })
}





func (r *RabbitMQ) setupExchangesAndQueue() error {
	_ , err:= r.Channel.QueueDeclare(
		"hello",  // name
	 true, //durable
	  false, //declare when use
	   false,  //exclusive
	   false,  // no-wait
	   nil,    //argument
	    ) 
if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (r *RabbitMQ) Close() {
	if r.conn != nil {
		r.conn.Close()
	}
	if r.Channel != nil {
		r.Channel.Close()
	}
}
