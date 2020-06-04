package middleware

import (
	"fmt"
	"net/http"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"os"
	"codproject/server/models"
	"encoding/json"
)

func ConnectDB() *sql.DB {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to DB")
	return db
}

func Signup(w http.ResponseWriter, r *http.Request) {
	// open database
	db := ConnectDB()
	creds := &models.Credentials{}

	// decode request body and check for error
	if err := json.NewDecoder(r.Body).Decode(creds); err != nil {
		fmt.Println(err)
		// return a 400 status
		w.WriteHeader(http.StatusBadRequest)
	}

	// insert into database and check for error
	if _, err := db.Query("insert into users (username, password) values ($1, $2)",
						creds.Username,
						creds.Password); err != nil {
							fmt.Println(err)
							// return 500 status
							w.WriteHeader(http.StatusInternalServerError)
						}
	fmt.Println("Signup successful")
}

func Login(w http.ResponseWriter, r *http.Request) {
	// open database
	db := ConnectDB()
	creds := &models.Credentials{}

	// decode request body and check for error
	if err := json.NewDecoder(r.Body).Decode(creds); err != nil {
		panic(err)
	}

	// query database for the password with username in request body
	storedCreds := &models.Credentials{}
	passFromDB := db.QueryRow("SELECT password FROM users WHERE username = $1", creds.Username).Scan(&storedCreds.Password)
	switch {
		case passFromDB == sql.ErrNoRows:
			fmt.Println("No user with the username: %d\n", creds.Username)
			// return 401 status
			w.WriteHeader(http.StatusUnauthorized)
		case passFromDB != nil:
			fmt.Println("Query error: %v\n", passFromDB)
			// return 500 status
			w.WriteHeader(http.StatusInternalServerError)
		default:
			if storedCreds.Password != creds.Password {
				fmt.Println("Access Denied: Wrong Password.")
				// return 401 status
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				fmt.Println("Access Granted.")
			}
	}
}


func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to home page")
	fmt.Println("Endpoint hit: HomePage")
}
