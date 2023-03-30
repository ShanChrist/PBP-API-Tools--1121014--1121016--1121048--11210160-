package controllers

import (
	"encoding/json"
	"net/http"
)

func userResponse(w http.ResponseWriter, users []User) {
	var response UserResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = users
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func SendRespond(w http.ResponseWriter, r *http.Request, message string, req interface{}) {
	var response Response
	response.Status = 200
	response.Message = message
	response.Data = req
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func SendRespondDoang(w http.ResponseWriter, r *http.Request, message string) {
	var response ResponseDoang
	response.Status = 200
	response.Message = message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func SendErrorResponse(w http.ResponseWriter, r *http.Request, message string) {
	var response ErrorResponse
	response.Status = 400
	response.Message = message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(response)
}
