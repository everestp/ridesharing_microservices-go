package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"ride-sharing/services/trip-service/internal/infrastructure/grpc"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/shared/env"
	"ride-sharing/shared/messaging"

	amqp "github.com/rabbitmq/amqp091-go"

	grpcsercver "google.golang.org/grpc"
)
var GrcpAddr = ":9093"

func main(){
	
rabbitMqURI :=env.GetString("RABBITMQ_URI", "ampq://guest:guest@localhost:5672/")
 
	inmemRepo :=repository.NewInmemRepository()
	svc :=service.NewService(inmemRepo)
	
     ctx ,cancel := context.WithCancel(context.Background())
	 defer cancel()
	 go func ()  {
		sigCh :=make(chan os.Signal ,1)
		  signal.Notify(sigCh,os.Interrupt, syscall.SIGTERM)
		  <-sigCh 
		  cancel()
	 }()
   lis ,err := net.Listen("tcp", GrcpAddr)
   if err != nil{
	log.Fatalf("Failed to listen : %v", err)
   }

   //RabbitMQ connection
    rabbitmq ,err :=messaging.NewRabbitMQ(rabbitMqURI)
	if err != nil{
		log.Fatal(err)
	}
	defer rabbitmq.Close()

	log.Println("Starting the Rabbit MQ Connection")
   
   //TODO intiliaze our grpc handler implementation

   //Starting the gRPC server
	 grpcserver  := grpcsercver.NewServer()
  
  grpc.NewGRPCHandler(grpcserver, svc)

  
  log.Printf("Starting gRPC server Trip service on port %s", lis.Addr().String())
  go func() {
if err := grpcserver.Serve(lis); err != nil{
	log.Printf("Failed to server : %v", err)
	cancel()
}
  }()
	// wait for the shutdown signal
   <-ctx.Done()
   log.Println("Shutting down grpc Server")
   grpcserver.GracefulStop()
}