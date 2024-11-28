package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

type TimeResponse struct {
	CurrentTime string `json:"current_time"`
}

type LogEntry struct {
	ID        int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	// Get database connection details from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Database connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Retry logic for database connection
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		var err error
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("Attempt %d: Error opening database: %v", i+1, err)
		} else {
			err = db.Ping()
			if err == nil {
				log.Printf("Successfully connected to database.")
				break // Successfully connected
			}
			log.Printf("Attempt %d: Error connecting to database: %v", i+1, err)
		}
		time.Sleep(5 * time.Second)
	}

	// Check if the connection was successful
	if db == nil {
		log.Fatal("Failed to connect to database after multiple attempts")
	}
	defer db.Close()

	// Set up the router
	router := mux.NewRouter()
	router.HandleFunc("/current-time", currentTimeHandler).Methods("GET")
	router.HandleFunc("/logs", getLogsHandler).Methods("GET")
	router.HandleFunc("/", indexHandler).Methods("GET")

	// Start the server
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func currentTimeHandler(w http.ResponseWriter, r *http.Request) {
	// Get Toronto time
	loc, err := time.LoadLocation("America/Toronto")
	if err != nil {
		http.Error(w, "Error loading time zone", http.StatusInternalServerError)
		log.Printf("Time zone error: %v", err)
		return
	}

	currentTime := time.Now().In(loc)

	// Insert the current time into the database
	_, err = db.Exec("INSERT INTO time_log (timestamp) VALUES (?)", currentTime)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("Database error: %v", err)
		return
	}

	// Prepare the response
	response := TimeResponse{
		CurrentTime: currentTime.Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getLogsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, timestamp FROM time_log ORDER BY id DESC")
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		log.Printf("Database query error: %v", err)
		return
	}
	defer rows.Close()

	var logs []LogEntry
	for rows.Next() {
		var entry LogEntry
		err := rows.Scan(&entry.ID, &entry.Timestamp)
		if err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			log.Printf("Row scan error: %v", err)
			return
		}
		logs = append(logs, entry)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Row iteration error", http.StatusInternalServerError)
		log.Printf("Row iteration error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
