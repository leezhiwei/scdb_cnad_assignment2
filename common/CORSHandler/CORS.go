package CORShandler

import (
	"net/http"
)

func SetCORSHeaders(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost") // Replace with your actual client origin
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // cors
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
