package main

import (
	"context"
	"log"
	

	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"time"
)


func main(){
	ctx := context.Background()
	inmemRepo :=repository.NewInmemRepository()
	svc :=service.NewService(inmemRepo)
	// mux := http.NewServeMux()
	// httphandler := h.HttpHandler{Service :svc}


fare :=&domain.RideFareModel{
	UserID: "42",
}
	 t,err := svc.CreateTrip(ctx, fare )
	 if err !=nil{
		log.Println(err)
	 }
	 log.Println(t)
	 // this is the temprorary program running for now
	 for {
		time.Sleep(time.Second)
       time.Sleep(time.Second)
	 }

}