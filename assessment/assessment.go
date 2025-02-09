package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

	CORShandler "github.com/leezhiwei/common/CORSHandler"
	"github.com/leezhiwei/common/mainhandler"
	"github.com/leezhiwei/common/ping"

	"github.com/gorilla/mux"
)

type Assessment struct {
	SeniorID       int  `json:"senior_id"`
	LegStrength    int  `json:"leg_strength`
	Vision         int  `json:"vision`
	Balance        int  `json:"balance`
	Medication     bool `json:"medication`
	HistoryOfFalls bool `json:"history_of_falls`
	KneeInjury     bool `json:"knee_injury`
}

type StringAssessment struct {
	SeniorID       string `json:"senior_id"`
	LegStrength    string `json:"leg_strength`
	Vision         string `json:"vision`
	Balance        string `json:"balance`
	Medication     string `json:"medication`
	HistoryOfFalls string `json:"history_of_falls`
	KneeInjury     string `json:"knee_injury`
}

var result string

// func main() {
// 	http.HandleFunc("/", homePage)
// 	http.HandleFunc("/submit", submitAssessment)
// 	http.ListenAndServe(":8080", nil)
// }

func homePage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("assessment.html")
	if err != nil {
		fmt.Println("Error loading template:", err)
		http.Error(w, "Error Loading Assessment Form. Please Try Again.", http.StatusInternalServerError)
		return
	}
	t.Execute(w, result)
}

func submitAssessment(w http.ResponseWriter, r *http.Request) {
	var preflight bool = CORShandler.SetCORSHeaders(&w, r)
	if preflight {
		return
	}

	if r.Method == http.MethodPost {
		var err error
		var tempStrAssessment StringAssessment
		err = json.NewDecoder(r.Body).Decode(&tempStrAssessment)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		assessment := Assessment{
			SeniorID:       parseToInt(tempStrAssessment.SeniorID),
			LegStrength:    parseToInt(tempStrAssessment.LegStrength),
			Vision:         parseToInt(tempStrAssessment.Vision),
			Balance:        parseToInt(tempStrAssessment.Balance),
			Medication:     tempStrAssessment.Medication == "1",
			HistoryOfFalls: tempStrAssessment.HistoryOfFalls == "1",
			KneeInjury:     tempStrAssessment.KneeInjury == "1",
		}

		result = calculateRisk(w, r, assessment)
		// http.Redirect(w, r, "/", http.StatusSeeOther)

		// Return JSON response instead of redirecting
		response := map[string]string{"risk_level": result}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func parseToInt(value string) int {
	var result int
	fmt.Sscanf(value, "%d", &result)
	return result
}

func calculateRisk(w http.ResponseWriter, r *http.Request, assessment Assessment) string {

	var OvrWellBg string
	score := 0

	score += (6 - assessment.LegStrength)
	score += (6 - assessment.Vision)
	score += (6 - assessment.Balance)
	if assessment.Medication {
		score += 1
	}
	if assessment.HistoryOfFalls {
		score += 2
	}
	if assessment.KneeInjury {
		score += 2
	}

	if score <= 2 {
		//return "Low Risk"
		OvrWellBg = "Low Risk"
	} else if score <= 4 {
		//return "Moderate Risk"
		OvrWellBg = "Moderate Risk"
	} else {
		//return "High Risk - Please consult a healthcare professional."
		OvrWellBg = "High Risk - Please consult a healthcare professional."
	}

	// Insert assessment in DB
	// Execute the query using a prepared statement with placeholder
	insertAsstQuery := "INSERT INTO Assessment (Overall_Wellbeing, SeniorID) VALUES (?, ?)"
	// Query the database
	insertBillresults, err := db.Query(insertAsstQuery, OvrWellBg, assessment.SeniorID)
	if err != nil {
		log.Println("Database query error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return "Database Error"
	}
	// Close results object from db.Query when function completes
	defer insertBillresults.Close()

	return OvrWellBg
}

func rechealthGuide(assessment Assessment) {
	// Check for leg strength, Vision and balance
	// Add values into database
	if assessment.LegStrength < 3 || assessment.Balance < 3 {
		// Recommended to improve leg strength and balance with guided exercises
	}
	if assessment.Vision < 3 {
		// Regular vision check-ups recommended or eye exercises
	}
	if assessment.HistoryOfFalls == true {
		// Video on how to prevent falls
	}
}

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
	var prefix string = "/api/v1/assessment"
	CORShandler.DebugMode = config.DebugMode
	CORShandler.Hostname = config.Hostname
	router := mux.NewRouter()
	router.HandleFunc(fmt.Sprintf("%s/ping", prefix), ping.PingHandler).Methods("GET", "OPTIONS")
	router.HandleFunc(fmt.Sprintf("%s/", prefix), homePage).Methods("POST", "OPTIONS")
	router.HandleFunc(fmt.Sprintf("%s/submit", prefix), submitAssessment).Methods("POST", "OPTIONS")
	router.Use(mainhandler.LogReq)
	log.Println(fmt.Sprintf("Assesment running at port %d", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
