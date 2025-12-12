package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	
	grpcsercver "google.golang.org/grpc"
)
var GrcpAddr = ":9092"

func main(){
	
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

   service := NewService()
   //TODO intiliaze our grpc handler implementation

   //Starting the gRPC server
	 grpcserver  := grpcsercver.NewServer()
	 NewGrpcHandler(grpcserver, service)
  

  
  log.Printf("Starting gRPC server Driver service on port %s", lis.Addr().String())
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