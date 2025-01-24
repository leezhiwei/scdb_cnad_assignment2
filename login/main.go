package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"database/sql"
	"encoding/json"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/leezhiwei/common/mainhandler"
	"github.com/leezhiwei/common/ping"
)

// user struct
type Senior struct {
	SeniorID int    `json:"senior_id"`
	Phone    string `json:"phone"`
	Name     string `json:"name"`
}

// login and register
// handleLoginOrRegister handles user login or registration based on phone number.
func handleLoginOrRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost") // Replace with your actual client origin
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	var err error

	// Parse JSON request body
	var reqData struct {
		Phone   string `json:"phone"`
		SMScode string `json:"smscode"`
	}

	err = json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid JSON data")
		return
	}

	// Validate input
	if reqData.Phone == "" || reqData.SMScode == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Phone and SMS code are required")
		return
	}

	// Check if the phone number exists in the database
	var senior Senior
	query := `SELECT SeniorID, Phone_number FROM Senior WHERE Phone_number = ?`
	err = db.QueryRow(query, reqData.Phone).Scan(&senior.SeniorID, &senior.Phone)

	if err != nil {
		if err == sql.ErrNoRows {
			// Phone number not found; proceed with registration
			registerQuery := `
                INSERT INTO Senior (Phone_number) 
                VALUES (?)
            `

			result, err := db.Exec(registerQuery, reqData.Phone)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Error registering user")
				return
			}

			SeniorID, _ := result.LastInsertId()
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"message":   "User registered successfully",
				"senior_id": fmt.Sprintf("%d", SeniorID),
			})
			return
		}

		// Other errors
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error checking user existence")
		return
	}

	// Phone number exists; proceed with login
	// Set cookie with UserID and expire after 24 hours
	http.SetCookie(w, &http.Cookie{
		Name:     "senior_id",
		Value:    fmt.Sprintf("%d", senior.SeniorID),
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":   "Login successful",
		"senior_id": fmt.Sprintf("%d", senior.SeniorID),
	})
}

// updateUserProfile allows users to update their personal details.
func updateUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost") // Replace with your actual client origin
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	var err error

	// Retrieve the user's ID from cookies
	userIDCookie, err := r.Cookie("user_id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unauthorized: Please log in")
		return
	}
	//store user id
	userID := userIDCookie.Value
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid user ID")
		return
	}

	// Parse JSON request body
	var reqData struct {
		Email string `json:"email"`
		Phone string `json:"phone"`
	}

	err = json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid JSON data")
		return
	}

	// Validate inputs
	if reqData.Email == "" && reqData.Phone == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "At least one field (email or phone) must be provided")
		return
	}

	// Update the user's profile in the database for user email and phone
	query := `
        UPDATE users
        SET Email = COALESCE(NULLIF(?, ''), Email),
            Phone = COALESCE(NULLIF(?, ''), Phone),
            UpdatedAt = NOW()
        WHERE UserID = ?
    `
	_, err = db.Exec(query, reqData.Email, reqData.Phone, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error updating user profile")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Profile updated successfully")
}

// viewUserProfile allows users to view their membership status and rental history.
func viewUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost") // Replace with your actual client origin
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	var err error
	// Retrieve the user's ID from cookies, not allow access if there is no cookies
	userIDCookie, err := r.Cookie("user_id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unauthorized: Please log in")
		return
	}
	//store user id
	userID := userIDCookie.Value
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid user ID")
		return
	}

	// Query user details
	var user User
	query := `
        SELECT Email, Phone, MembershipTierID, MembershipPoint, CreatedAt, UpdatedAt
        FROM users
        WHERE UserID = ?
    `
	err = db.QueryRow(query, userID).Scan(&user.Email, &user.Phone, &user.MembershipTierID, &user.MembershipPoint, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error retrieving user profile")
		return
	}

	// Rental struct user for retrice rental history
	type Reservation struct {
		ReservationID int       `json:"reservation_id"`
		VehicleID     int       `json:"vehicle_id"`
		StartTime     time.Time `json:"start_time"`
		EndTime       time.Time `json:"end_time"`
		Status        string    `json:"status"`
	}
	//retrice from database according to user ID
	reservations := []Reservation{}
	query = `
        SELECT ReservationID, VehicleID, StartTime, EndTime, Status
        FROM reservations
        WHERE UserID = ?
        ORDER BY CreatedAt DESC
    `
	rows, err := db.Query(query, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error retrieving rental history")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var reservation Reservation
		err := rows.Scan(&reservation.ReservationID, &reservation.VehicleID, &reservation.StartTime, &reservation.EndTime, &reservation.Status)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error processing rental history")
			return
		}
		reservations = append(reservations, reservation)
	}

	// Create response
	response := map[string]interface{}{
		"user":          user,
		"rentalHistory": reservations,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Password string `json:"password"`
		User     string `json:"user"`
		Port     int    `json:"port"`
		DBName   string `json:"dbname"`
	} `json:"database"`
	ServPort int `json:"server_port"`
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
			Port: 5432,
		},
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
func main() {
	var errdb error
	var config Config = GetConfig()
	var connstring = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.Database.User, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.DBName)
	// connection string
	db, errdb = sql.Open("mysql", connstring) // make sql connection
	if errdb != nil {                         // if error with db
		log.Fatal("Unable to connect to database, error: ", errdb) // print err
	}
	var port int = config.ServPort
	var prefix string = "/api/v1/login"
	router := mux.NewRouter()
	router.HandleFunc(fmt.Sprintf("%s/ping", prefix), ping.PingHandler).Methods("GET")
	router.Use(mainhandler.LogReq)
	log.Println(fmt.Sprintf("Login Server running at port %d", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
