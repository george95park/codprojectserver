package middleware

import (
	"fmt"
	"time"
	"net/http"
	"database/sql"
	"encoding/json"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"codproject/server/models"
	"codproject/server/config"
)

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
	db := config.ConnectDB()
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
	// create new uuid for cookie
	sessionToken := uuid.Must(uuid.NewV4()).String()

	// insert into credentials table and check for error
	_,err = db.Query("insert into credentials (username, password, token) values ($1, $2, $3)",creds.Username,string(hashedPassword),sessionToken)
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
	res := models.Response{}
	err = db.QueryRow("select user_id from credentials where username = $1", creds.Username).Scan(&res.User_Id)
	if err != nil {
		panic(err)
	}
	res.Message = "success"
	res.Username = creds.Username
	res.Logged_In = true
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
	db := config.ConnectDB()
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
	err := db.QueryRow("select password,user_id from credentials where username = $1", creds.Username).Scan(&storedCreds.Password, &storedCreds.User_Id)
	switch {
	case err == sql.ErrNoRows:
			fmt.Println("No user with the username: %d\n", creds.Username)
			// return 401 status
			w.WriteHeader(http.StatusUnauthorized)
			return
		case err != nil:
			fmt.Println("Query error: %v\n", err)
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
				_, err := db.Query("update credentials set token=$1 where username=$2", sessionToken, creds.Username)
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
					User_Id: storedCreds.User_Id,
					Logged_In: true,
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
	db := config.ConnectDB()
	defer db.Close()
	creds := &models.Credentials{}
	// decode request body and check for error
	if err := json.NewDecoder(r.Body).Decode(creds); err != nil {
		fmt.Println(err)
		// return a 400 status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res := models.Response{}
	// get the username related to token
	userFromDB := db.QueryRow("select user_id,username from credentials where token=$1", creds.Token).Scan(&res.User_Id, &res.Username)
	switch {
		// if there were no user with this token
		case userFromDB == sql.ErrNoRows:
			fmt.Println("No user with the token: %d\n", creds.Token)
			// return 401 status
			w.WriteHeader(http.StatusUnauthorized)
			return
		case userFromDB != nil:
			fmt.Println("Query error: %v\n", userFromDB)
			// return 500 status
			w.WriteHeader(http.StatusInternalServerError)
			return
		default:
			res.Logged_In = true
			res.Message = "success"
			json.NewEncoder(w).Encode(res)
	}
}
