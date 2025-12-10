package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	grpcclients "ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
	"time"
)

func handleTripPreview(w http.ResponseWriter , r *http.Request){
     time.Sleep(time.Second * 9)
	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err !=nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	// validation
	if reqBody.UserID == ""{
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return 
	}
	jsonBody ,_:= json.Marshal(reqBody)
	reader := bytes.NewReader(jsonBody)

// grpc -> NOT GOOD IDEA TO PUT gRPC connection here as traffic grow it affect other service
tripService ,err := grpcclients.NewTripServiceClient()
if err != nil {
	log.Fatal(err)
}
defer tripService.Close()
// tripService.Client.PreviewTrip()



	// TODO : Call Trip Service
	resp ,err := http.Post("https://trip-service:8083/preview", "applications/json",reader)
	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()
	var respbody any
	if err := json.NewDecoder(resp.Body).Decode(&respbody); err != nil{
		http.Error(w, "failed to parse json from trip service",http.StatusOK)
	}
	response :=contracts.APIResponse{Data: respbody}
	writeJSON(w, http.StatusCreated, response)
	//

}