package main

import "net/http"

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write a simple JSON response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))

}
