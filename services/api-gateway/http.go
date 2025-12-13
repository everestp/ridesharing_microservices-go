package main

import (
	"encoding/json"
	"log"
	"net/http"
	grpcclients "ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
)


func handleTripStart(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "handleTripStart")
	defer span.End()

	var reqBody startTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// Why we need to create a new client for each connection:
	// because if a service is down, we don't want to block the whole application
	// so we create a new client for each connection
	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	// Don't forget to close the client to avoid resource leaks!
	defer tripService.Close()

	trip, err := tripService.Client.CreateTrip(ctx, reqBody.toProto())
	if err != nil {
		log.Printf("Failed to start a trip: %v", err)
		http.Error(w, "Failed to start trip", http.StatusInternalServerError)
		return
	}

	response := contracts.APIResponse{Data: trip}

	writeJSON(w, http.StatusCreated, response)
}

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
