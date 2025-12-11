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

	grpcsercver "google.golang.org/grpc"
)
var GrcpAddr = ":9093"

func main(){
	


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