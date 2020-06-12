package middleware

import (
	"os"
	"fmt"
	"time"
	"net/http"
	"database/sql"
	"encoding/json"
	"github.com/lib/pq"
	"github.com/joho/godotenv"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"codproject/server/models"
)

// opens access to database
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

// creates loadout in database
func CreateLoadout(w http.ResponseWriter, r *http.Request) {
	db := ConnectDB()
	defer db.Close()
	load := &models.Loadout{}
	// Unmarshall request body
	if err := json.NewDecoder(r.Body).Decode(load); err != nil {
		panic(err)
	}

	// insert new loadout to database
	if _, err := db.Query("insert into loadouts (username, gunname, attachments, description) values ($1, $2, $3, $4)",
		load.Username,
		load.Gunname,
		pq.Array(load.Attachments),
		load.Description); err != nil {
			panic(err)
		}
	fmt.Println("Loadout created")
	json.NewEncoder(w).Encode(load)
}

// gets loadouts from database according to user
func GetLoadouts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	db := ConnectDB()
	defer db.Close()
	currUser := mux.Vars(r)["user"]


	// Get the rows according to user
	rows,err := db.Query("select id,gunname,attachments,description from loadouts where username=$1",currUser)
	if err != nil {
		panic(err)
	}

	// For each row, create new loadout object and append to the result
	userLoadouts := []models.Loadout{}
	for rows.Next() {
		l := models.Loadout{}
		l.Username = currUser
		err = rows.Scan(&l.Id, &l.Gunname, pq.Array(&l.Attachments), &l.Description)
		if err != nil {
			panic(err)
		}
		userLoadouts = append(userLoadouts, l)
	}
	fmt.Println("Received loadouts.")

	// return response
	json.NewEncoder(w).Encode(userLoadouts)
}

func DeleteLoadout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// open database
	db := ConnectDB()
	defer db.Close()
	load := &models.Loadout{}
	if err := json.NewDecoder(r.Body).Decode(load); err != nil {
		fmt.Println(err)
		// return a 400 status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res,err := db.Exec("delete from loadouts where id=$1", load.Id)
	if err != nil {
		fmt.Println(err)
		// return 500 status
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Println("Total rows affected: %v", rowsAffected)
	json.NewEncoder(w).Encode(rowsAffected)
}

func UpdateLoadout(w http.ResponseWriter, r *http.Request) {
	// open database
	db := ConnectDB()
	defer db.Close()
	load := &models.Loadout{}
	if err := json.NewDecoder(r.Body).Decode(load); err != nil {
		fmt.Println(err)
		// return a 400 status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res, err := db.Exec("update loadouts set username=$2, gunname=$3, attachments=$4, description=$5 where id=$1",
		load.Id,
		load.Username,
		load.Gunname,
		pq.Array(load.Attachments),
		load.Description)
	if err != nil {
		panic(err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Println("Total rows affected: %v", rowsAffected)
	json.NewEncoder(w).Encode(rowsAffected)
}

// Sign-up handler
func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// open database
	db := ConnectDB()
	defer db.Close()
	creds := &models.Credentials{}

	// decode request body and check for error
	if err := json.NewDecoder(r.Body).Decode(creds); err != nil {
		fmt.Println(err)
		// return a 400 status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
	if err != nil {
		panic(err)
	}

	// insert into credentials table and check for error
	_,err = db.Query("insert into credentials (username, password) values ($1, $2)",creds.Username,string(hashedPassword))
	if err != nil {
		fmt.Println(err)
		// return 500 status
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// insert into USERS table and check for error
	if _, err := db.Query("insert into users (username) values ($1)", creds.Username); err != nil {
		fmt.Println(err)
		// return 500 status
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// create new uuid for cookie
	sessionToken := uuid.Must(uuid.NewV4()).String()
	// insert into tokens table to save user's session token
	_,err = db.Query("insert into tokens (token, username) values ($1, $2)", sessionToken, creds.Username)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	http.SetCookie(w, &http.Cookie {
		Name: "session_token",
		Value: sessionToken,
		Expires: time.Now().Add(365 * 24 * time.Hour),
		Domain: "/",
	})
	res := models.Response {
		Message: "success",
		Username: creds.Username,
		LoggedIn: true,
	}
	json.NewEncoder(w).Encode(res)
}
// Login handler
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	// handles preflight request before the actual request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// open database
	db := ConnectDB()
	defer db.Close()
	creds := &models.Credentials{}
	// decode request body and check for error
	if err := json.NewDecoder(r.Body).Decode(creds); err != nil {
		fmt.Println(err)
		// return a 400 status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// query database for the password with username in request body
	storedCreds := &models.Credentials{}
	passFromDB := db.QueryRow("select password from credentials where username = $1", creds.Username).Scan(&storedCreds.Password)
	switch {
		case passFromDB == sql.ErrNoRows:
			fmt.Println("No user with the username: %d\n", creds.Username)
			// return 401 status
			w.WriteHeader(http.StatusUnauthorized)
			return
		case passFromDB != nil:
			fmt.Println("Query error: %v\n", passFromDB)
			// return 500 status
			w.WriteHeader(http.StatusInternalServerError)
			return
		default:
			if err := bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
				fmt.Println("Access Denied: Wrong Password.")
				// return 401 status
				w.WriteHeader(http.StatusUnauthorized)
				return
			} else {
				fmt.Println("Access granted.")
				// update session token in tokens table
				sessionToken := uuid.Must(uuid.NewV4()).String()
				_, err := db.Query("update tokens set token=$1 where username=$2", sessionToken, creds.Username)
				if err != nil {
					fmt.Println(err)

					// return 500 status
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				http.SetCookie(w, &http.Cookie {
					Name: "session_token",
					Value: sessionToken,
					Expires: time.Now().Add(365 * 24 * time.Hour),
				})
				res := models.Response {
					Message: "success",
					Username: creds.Username,
					LoggedIn: true,
				}
				json.NewEncoder(w).Encode(res)
			}
	}
}

func GetSessionTokenUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// handles preflight request before the actual request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	// open database
	db := ConnectDB()
	defer db.Close()
	token := &models.Token{}
	// decode request body and check for error
	if err := json.NewDecoder(r.Body).Decode(token); err != nil {
		fmt.Println(err)
		// return a 400 status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// get the username related to token
	res := models.Response{}
	userFromDB := db.QueryRow("select username from tokens where token=$1", token.Token).Scan(&res.Username)
	switch {
		// if there were no user with this token
		case userFromDB == sql.ErrNoRows:
			fmt.Println("No user with the token: %d\n", token.Token)
			res.Message = "not found"
			res.LoggedIn = false
			json.NewEncoder(w).Encode(res)
		case userFromDB != nil:
			fmt.Println("Query error: %v\n", userFromDB)
			// return 500 status
			w.WriteHeader(http.StatusInternalServerError)
			return
		default:
			res.Message = "found"
			res.LoggedIn = true
			json.NewEncoder(w).Encode(res)
	}
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to home page")
	fmt.Println("Endpoint hit: HomePage")
}
