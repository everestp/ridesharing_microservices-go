package main

import (
	"encoding/json"
	"net/http"
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
	// TODO : Call Trip Service
	response :=contracts.APIResponse{Data: "ok"}
	writeJSON(w, http.StatusCreated, response)
	//

}