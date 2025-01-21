package mainhandler

import (
	"fmt"
	"log"
	"net/http"
)

func LogReq(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(fmt.Sprintf("From %s IP, accessing %s", r.RemoteAddr, r.URL))
		next.ServeHTTP(w, r)
	})
}
