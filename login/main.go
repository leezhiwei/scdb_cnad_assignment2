package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leezhiwei/common/mainhandler"
	"github.com/leezhiwei/common/ping"
)

func main() {
	var port int = 8080
	router := mux.NewRouter()
	router.HandleFunc("/ping", ping.PingHandler).Methods("GET")
	router.Use(mainhandler.LogReq)
	log.Println(fmt.Sprintf("Login Server running at port %d", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
