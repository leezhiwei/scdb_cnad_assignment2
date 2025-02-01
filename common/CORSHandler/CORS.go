package CORShandler

import (
	"net/http"
)

func SetCORSHeaders(w *http.ResponseWriter, r *http.Request) bool { // if fn returns true, cors preflight
	(*w).Header().Set("Access-Control-Allow-Origin", "*") // Replace with your actual client origin
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // cors
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	if r.Method == "OPTIONS" {
		(*w).WriteHeader(http.StatusOK)
		return true
	}
	return false
}
