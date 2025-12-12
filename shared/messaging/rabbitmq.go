package messaging

import (
	"fmt"
	

	ampq "github.com/rabbitmq/amqp091-go"
)


type RabbitMQ struct {
	conn *ampq.Connection
}

func NewRabbitMQ(uri string)(*RabbitMQ , error){
	 //RabbitMQ connection
    conn , err :=ampq.Dial(uri)
	if err != nil {
		return  nil , fmt.Errorf("Failed to Connect to Rabbit mq : %v", err)
		
	}
	rmq :=&RabbitMQ{
		conn: conn,
	}
	return  rmq ,nil
}


func (r *RabbitMQ) Close(){
	if r.conn !=nil{
		r.conn.Close()
	}
}
