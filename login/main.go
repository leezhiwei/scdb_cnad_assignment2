package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"database/sql"
	"encoding/json"
	"math/rand"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	CORShandler "github.com/leezhiwei/common/CORSHandler"
	"github.com/leezhiwei/common/mainhandler"
	"github.com/leezhiwei/common/ping"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

// user struct
type Senior struct {
	SeniorID int    `json:"senior_id"`
	Phone    string `json:"phone"`
	Name     string `json:"name"`
}

var (
	otpStore = make(map[string]string)
	mutex    = &sync.Mutex{}
)

func generateRandomNumber() string {
	if config.DebugMode {
		return "123456" // Fixed number for debugging
	}
	time.Now().UnixNano()
	return fmt.Sprintf("%06d", rand.Intn(999999-1)+1)
}

func sendSMS(to string, code string) error {
	var accountSid string = config.APIToken.Twilio.AccountSID
	var authToken string = config.APIToken.Twilio.AuthToken
	var from string = config.APIToken.Twilio.FromNumber

	// Create a Twilio client
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
	to = "+65 " + to // add SG prefix
	// Create the message
	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(from)
	params.SetBody(fmt.Sprintf("Your login code for SDCB is: %s", code))

	// Send the SMS
	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	log.Printf("SMS sent successfully: SID=%s", *resp.Sid)
	return nil
}

func handleSMS(w http.ResponseWriter, r *http.Request) {
	var preflight bool = CORShandler.SetCORSHeaders(&w, r)
	if preflight {
		return
	}

	var req struct {
		Phone string `json:"phone"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Phone == "" {
		http.Error(w, "Phone number required", http.StatusBadRequest)
		return
	}

	otp := generateRandomNumber()
	mutex.Lock()
	otpStore[req.Phone] = otp
	mutex.Unlock()

	if !config.DebugMode {
		if err := sendSMS(req.Phone, otp); err != nil {
			http.Error(w, "Failed to send SMS", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "OTP sent"})
}

// login and register
// handleLoginOrRegister handles user login or registration based on phone number.
func handleLogin(w http.ResponseWriter, r *http.Request) {
	var preflight bool = CORShandler.SetCORSHeaders(&w, r)
	if preflight {
		return
	}

	var req struct {
		Phone   string `json:"phone"`
		SMScode string `json:"smscode"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	storedCode, exists := otpStore[req.Phone]
	mutex.Unlock()

	fmt.Println(exists)
	if !exists || storedCode != req.SMScode {
		http.Error(w, "Invalid OTP", http.StatusUnauthorized)
		return
	}

	var senior Senior
	//senior.SeniorID = 1
	query := `SELECT SeniorID, Phone_number FROM Senior WHERE Phone_number = ?`
	err := db.QueryRow(query, req.Phone).Scan(&senior.SeniorID, &senior.Phone)
	// If user doesn't exist, register automatically
	if err == sql.ErrNoRows {
		// Insert new user
		insertQuery := `INSERT INTO Senior (Phone_number) VALUES (?)`
		result, err := db.Exec(insertQuery, req.Phone)
		if err != nil {
			http.Error(w, "Error registering user", http.StatusInternalServerError)
			return
		}

		// Get new user ID
		newID, err := result.LastInsertId()
		if err != nil {
			http.Error(w, "Error retrieving user ID", http.StatusInternalServerError)
			return
		}
		senior.SeniorID = int(newID)
		senior.Phone = req.Phone
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "senior_id",
		Value:    fmt.Sprintf("%d", senior.SeniorID),
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":   "Login successful",
		"senior_id": fmt.Sprintf("%d", senior.SeniorID),
	})
}

// emergency contact struct
type EmergencyContact struct {
	Contactid      int    `json:"contact_id"`
	ContactName    string `json:"contactname"`
	ContactNumbert string `json:"contactnumber"`
	SeniorID       int    `json:"senior_id"`
}

// In-memory storage for emergency contacts
var contacts = make(map[string]EmergencyContact)
var mu sync.Mutex

func AddEmergencyContact(w http.ResponseWriter, r *http.Request) {
	var preflight bool = CORShandler.SetCORSHeaders(&w, r)
	if preflight {
		return
	}

	var req struct {
		ContactName   string `json:"contactname"`
		ContactNumber string `json:"contactnumber"`
		SeniorID      int    `json:"senior_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.ContactName == "" || req.ContactNumber == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	seniorID := req.SeniorID
	query := `INSERT INTO Emergency_Contact (ContactName, ContactNumber, SeniorID) VALUES (?, ?, ?)`
	result, err := db.Exec(query, req.ContactName, req.ContactNumber, seniorID)
	if err != nil {
		http.Error(w, "Error inserting emergency contact", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	newID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Error retrieving inserted ID", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":             "Emergency contact added successfully",
		"emergencycontact_id": newID,
	})
}

func ListEmergencyContact(w http.ResponseWriter, r *http.Request) {
	var preflight bool = CORShandler.SetCORSHeaders(&w, r)
	if preflight {
		return
	}

	var req struct {
		SeniorID int `json:"senior_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.SeniorID == 0 {
		http.Error(w, "Senior ID is required", http.StatusBadRequest)
		return
	}

	seniorID := req.SeniorID

	//Query the database for emergency contacts
	query := `SELECT Contactid, ContactName, ContactNumber, SeniorID FROM Emergency_Contact WHERE SeniorID = ?`
	rows, err := db.Query(query, seniorID)
	fmt.Println(seniorID)
	if err != nil {
		http.Error(w, "Error retrieving emergency contacts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Collect results
	var contacts []EmergencyContact
	for rows.Next() {
		var contact EmergencyContact
		if err := rows.Scan(&contact.Contactid, &contact.ContactName, &contact.ContactNumbert, &contact.SeniorID); err != nil {
			http.Error(w, "Error scanning database results", http.StatusInternalServerError)
			return
		}
		contacts = append(contacts, contact)
	}
	fmt.Println("Contacts found:", contacts)
	// Return the results as JSON
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contacts)
}

func DeleteEmergencyContact(w http.ResponseWriter, r *http.Request) {
	var preflight bool = CORShandler.SetCORSHeaders(&w, r)
	if preflight {
		return
	}

	var req struct {
		ContactID int `json:"contact_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.ContactID == 0 {
		http.Error(w, "Contact ID is required", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM Emergency_Contact WHERE Contactid = ?`
	result, err := db.Exec(query, req.ContactID)
	if err != nil {
		log.Println("Error executing DELETE query:", err)
		http.Error(w, "Error deleting emergency contact", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error retrieving affected rows:", err)
		http.Error(w, "Error retrieving affected rows", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		log.Println("No contact found with the given ID:", req.ContactID)
		http.Error(w, "No contact found with the given ID", http.StatusNotFound)
		return
	}

	log.Println("Successfully deleted contact:", req.ContactID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Emergency contact deleted successfully"})

}

func getSenior(w http.ResponseWriter, r *http.Request) {
	var preflight bool = CORShandler.SetCORSHeaders(&w, r)
	if preflight {
		return
	}
	var req struct {
		SeniorID int `json:"senior_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var seniordata Senior

	seniordata.SeniorID = req.SeniorID

	query := `SELECT Phone_number, Name FROM Senior WHERE SeniorID = ?`
	err := db.QueryRow(query, req.SeniorID).Scan(&seniordata.Phone, &seniordata.Name)
	// unable to find senior
	if err == sql.ErrNoRows {
		log.Println("Unable to find SeniorID in database.")
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Unable to find Senior",
		})
		return
	}
	// can find
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(seniordata)
	return
}

func updateSenior(w http.ResponseWriter, r *http.Request) {
	var err error
	var preflight bool = CORShandler.SetCORSHeaders(&w, r)
	if preflight {
		return
	}

	var seniordata Senior
	if err := json.NewDecoder(r.Body).Decode(&seniordata); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if seniordata.Name == "" {
		http.Error(w, "Please update your name.", http.StatusBadRequest)
		return
	}

	query := `
        UPDATE Senior
        SET Phone_number = COALESCE(NULLIF(?, ''), Phone_number),
            Name = COALESCE(NULLIF(?, ''), Name)
        WHERE SeniorID = ?
    `
	_, err = db.Exec(query, seniordata.Phone, seniordata.Name, seniordata.SeniorID) // run

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error updating Senior profile") // if err, error
		return
	}

	w.WriteHeader(http.StatusOK) // else success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Updated Senior",
	})
	return
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
	// connection string
	db, errdb = sql.Open("mysql", connstring) // make sql connection
	if errdb != nil {                         // if error with db
		log.Fatal("Unable to connect to database, error: ", errdb) // print err
	}
	var port int = config.ServPort
	var prefix string = "/api/v1/login"
	CORShandler.DebugMode = config.DebugMode
	CORShandler.Hostname = config.Hostname
	router := mux.NewRouter()
	router.HandleFunc(fmt.Sprintf("%s/ping", prefix), ping.PingHandler).Methods("GET", "OPTIONS")
	router.HandleFunc(fmt.Sprintf("%s/sendsms", prefix), handleSMS).Methods("POST", "OPTIONS")
	router.HandleFunc(fmt.Sprintf("%s/login", prefix), handleLogin).Methods("POST", "OPTIONS")
	router.HandleFunc(fmt.Sprintf("%s/addemergencycontact", prefix), AddEmergencyContact).Methods("POST", "OPTIONS")
	router.HandleFunc(fmt.Sprintf("%s/listemergencycontact", prefix), ListEmergencyContact).Methods("POST", "OPTIONS")
	router.HandleFunc(fmt.Sprintf("%s/deleteemergencycontact", prefix), DeleteEmergencyContact).Methods("POST", "OPTIONS")
	router.HandleFunc(fmt.Sprintf("%s/getsenior", prefix), getSenior).Methods("POST", "OPTIONS")
	router.HandleFunc(fmt.Sprintf("%s/updatesenior", prefix), updateSenior).Methods("POST", "OPTIONS")
	router.Use(mainhandler.LogReq)
	log.Println(fmt.Sprintf("Login Server running at port %d", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
