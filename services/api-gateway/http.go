package main

import (
	"encoding/json"
	"log"
	"net/http"
	grpcclients "ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {

	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		// FIXED: return proper JSON instead of plain text
		WriteJSONError(w, http.StatusBadRequest, "failed to parse JSON data") // FIXED
		return
	}
	defer r.Body.Close()

	// validation
	if reqBody.UserID == "" {
		// FIXED: return proper JSON error
		WriteJSONError(w, http.StatusBadRequest, "User ID is required") // FIXED
		return
	}

	// grpc client initialization
	tripService, err := grpcclients.NewTripServiceClient()
	if err != nil {
		// FIXED: don't use log.Fatal, return error JSON instead
		log.Printf("Failed to connect to trip service: %v", err)        // FIXED
		WriteJSONError(w, http.StatusInternalServerError, "trip service unavailable") // FIXED
		return
	}
	defer tripService.Close()

	tripPreview, err := tripService.Client.PreviewTrip(r.Context(), reqBody.toProto())
	if err != nil {
		// FIXED: JSON error + ensured return
		log.Printf("Failed to preview a trip: %v", err)                   // FIXED
		WriteJSONError(w, http.StatusInternalServerError, "Failed to preview Trip") // FIXED
		return
	}

	// SUCCESS RESPONSE - unchanged
	response := contracts.APIResponse{Data: tripPreview}
	writeJSON(w, http.StatusCreated, response)
}
