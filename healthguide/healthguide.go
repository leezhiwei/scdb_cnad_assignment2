package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	CORShandler "github.com/leezhiwei/common/CORSHandler"
	"github.com/leezhiwei/common/mainhandler"
	"github.com/leezhiwei/common/ping"

	"github.com/gorilla/mux"
)

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Password string `json:"password"`
		User     string `json:"user"`
		Port     int    `json:"port"`
		DBName   string `json:"dbname"`
	} `json:"database"`
	APIToken struct {
		Twilio struct {
			FromNumber string `json:"fromNumber"`
			AccountSID string `json:"accountSID"`
			AuthToken  string `json:"authToken"`
		} `json:"twilio"`
	} `json:"api_tokens"`
	Hostname  string `json:"hostname"`
	DebugMode bool   `json:"debugMode"`
	ServPort  int    `json:"server_port"`
}

func GetConfig() Config {
	config := Config{
		Database: struct {
			Host     string `json:"host"`
			Password string `json:"password"`
			User     string `json:"user"`
			Port     int    `json:"port"`
			DBName   string `json:"dbname"`
		}{
			Port: 3306,
		},
		Hostname:  "localhost",
		DebugMode: true,
		ServPort:  8080,
	}

	configFile, err := os.Open("./config.json")
	defer configFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

var db *sql.DB // global var
var config Config

// Health Guide Structure
type healthguide struct {
	SeniorID               int    `json:"senior_id"`
	HealthGuideID          int    `json:"healthguide_id"`
	HealthGuideDescription string `json:"healthguide_description"`
	HealthGuideVideoLink   string `json:"healthguide_videolink"`
	OverallWellbeing       string `json:"overall_wellbeing"`
}

func offerHealthGuide(w http.ResponseWriter, r *http.Request) {
	var preflight bool = CORShandler.SetCORSHeaders(&w, r)
	if preflight {
		return
	}

	//var tempHG healthguide

	// Extract SeniorID from query parameters
	seniorIDStr := r.URL.Query().Get("senior_id")
	if seniorIDStr == "" {
		http.Error(w, "Missing SeniorID", http.StatusBadRequest)
		return
	}

	// Convert SeniorID to integer
	seniorID, err := strconv.Atoi(seniorIDStr)
	if err != nil {
		http.Error(w, "Invalid SeniorID", http.StatusBadRequest)
		return
	}

	// Select health guides in DB
	// Select query using a prepared statement with placeholder
	query := `
        SELECT MAX(a.AssessmentID), hg.HealthGuideDescription, hg.HealthGuideVideoLink
		FROM HealthGuide hg
		INNER JOIN Assessment a ON hg.Overall_Wellbeing = a.Overall_Wellbeing
		WHERE a.SeniorID = ?;`

	// Execute the query
	rows, err := db.Query(query, seniorID)
	if err != nil {
		log.Println("Database query error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Slice to store health guides
	var healthGuides []healthguide

	// Add health guide description and video link in slice
	for rows.Next() {
		var guide healthguide
		if err := rows.Scan(&guide.HealthGuideDescription, &guide.HealthGuideVideoLink); err != nil {
			log.Println("Row scan error:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		healthGuides = append(healthGuides, guide)
	}

	// Check for errors after iterating
	if err := rows.Err(); err != nil {
		log.Println("Row iteration error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Encode results as JSON and return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(healthGuides)

}

func main() {
	var errdb error
	config = GetConfig()
	var connstring = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", config.Database.User, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.DBName)
	// Connection string
	db, errdb = sql.Open("mysql", connstring)
	// If error with db
	if errdb != nil {
		// Print err
		log.Fatal("Unable to connect to database, error: ", errdb)
	}
	var port int = config.ServPort
	var prefix string = "/api/v1/healthguide"
	CORShandler.DebugMode = config.DebugMode
	CORShandler.Hostname = config.Hostname
	router := mux.NewRouter()
	router.HandleFunc(fmt.Sprintf("%s/ping", prefix), ping.PingHandler).Methods("GET", "OPTIONS")
	router.HandleFunc(fmt.Sprintf("%s/suggestions", prefix), offerHealthGuide).Methods("GET", "OPTIONS")
	router.Use(mainhandler.LogReq)
	log.Println(fmt.Sprintf("Assesment running at port %d", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
