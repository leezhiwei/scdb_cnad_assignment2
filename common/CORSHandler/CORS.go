package CORShandler

import (
	"fmt"
	"net/http"
)

var DebugMode bool = true
var Hostname string = "localhost"

func SetCORSHeaders(w *http.ResponseWriter, r *http.Request) bool { // if fn returns true, cors preflight
	if DebugMode {
		(*w).Header().Set("Access-Control-Allow-Origin", "https://localhost")
	} else {
		(*w).Header().Set("Access-Control-Allow-Origin", fmt.Sprintf("https://%s", Hostname))
	}
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // cors
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	(*w).Header().Set("Access-Control-Max-Age", "1")
	if r.Method == "OPTIONS" {
		(*w).WriteHeader(http.StatusOK)
		return true
	}
	return false
}
