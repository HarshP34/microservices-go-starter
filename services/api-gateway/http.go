package main

import (
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {

	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if reqBody.UserID == "" {
		http.Error(w, "user ID is required", http.StatusBadRequest)
		return
	}

	// jsonBody, _ := json.Marshal(reqBody)
	// reader := bytes.NewReader(jsonBody)

	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	defer tripService.Close()

	tripServicePreview, err := tripService.Client.PreviewTrip(r.Context(), reqBody.toProto())
	if err != nil {
		log.Printf("Failed to preview a trip: %v", err)
		http.Error(w, "Failed to preview a trip", http.StatusInternalServerError)
		return
	}

	// resp, err:=http.Post("http://trip-service:8083/preview", "application/json", reader)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// defer resp.Body.Close()
	
	// var responseBody any
	// if err := json.NewDecoder(resp.Body).Decode(&responseBody);  err != nil {
	// 	http.Error(w, "Failed to decode response from trip service", http.StatusInternalServerError)
	// 	return
	// }

	response := contracts.APIResponse{Data: tripServicePreview}
	
	writeJSON(w, http.StatusCreated, response)
}


func handleTripStart(w http.ResponseWriter, r *http.Request) {
		var reqBody startTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if reqBody.UserID == "" {
		http.Error(w, "user ID is required", http.StatusBadRequest)
		return
	}
	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	defer tripService.Close()

	tripServiceStart, err := tripService.Client.CreateTrip(r.Context(), reqBody.toProto())
	if err != nil {
		log.Printf("Failed to start a trip: %v", err)
		http.Error(w, "Failed to start a trip", http.StatusInternalServerError)
		return
	}

	response := contracts.APIResponse{Data: tripServiceStart}
	writeJSON(w, http.StatusCreated, response)
}