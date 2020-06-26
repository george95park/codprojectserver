package middleware

import (
	"os"
	"fmt"
	"time"
	"net/http"
	"database/sql"
	"encoding/json"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"codproject/server/models"
	"codproject/server/config"
	"github.com/dgrijalva/jwt-go"
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

	// insert into credentials table and check for error
	_,err = db.Query("insert into credentials (username, password, token) values ($1, $2, $3)",creds.Username,string(hashedPassword),"")
	if err != nil {
		fmt.Println(err)
		// return 500 status
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// jwt token
	token,user_id := createToken(creds.Username, db)
	_,err = db.Exec("update credentials set token=$1 where user_id=$2", token, user_id)
	if err != nil {
		panic(err)
	}

	http.SetCookie(w, &http.Cookie {
		Name: "token",
		Value: token,
		Expires: time.Now().Add(365 * 24 * time.Hour),
		//HttpOnly: true,
	})
	res := models.Response{
		Message: "success",
		Username: creds.Username,
		User_Id: user_id,
		Logged_In: true,
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
	var storedPassword string
	err := db.QueryRow("select password from credentials where username = $1", creds.Username).Scan(&storedPassword)
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
			if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password)); err != nil {
				fmt.Println("Access Denied: Wrong Password.")
				// return 401 status
				w.WriteHeader(http.StatusUnauthorized)
				return
			} else {
				fmt.Println("Access granted.")
				token, user_id := createToken(creds.Username, db)
				_, err := db.Exec("update credentials set token=$1 where user_id=$2", token, user_id)
				if err != nil {
					fmt.Println(err)

					// return 500 status
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				http.SetCookie(w, &http.Cookie {
					Name: "token",
					Value: token,
					Expires: time.Now().Add(365 * 24 * time.Hour),
					//HttpOnly: true,
				})
				res := models.Response {
					Message: "success",
					Username: creds.Username,
					User_Id: user_id,
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
	err := db.QueryRow("select user_id,username from credentials where token=$1", creds.Token).Scan(&res.User_Id, &res.Username)
	switch {
		// if there were no user with this token
		case err == sql.ErrNoRows:
			fmt.Println("No user with the token: %d\n", creds.Token)
			// return 401 status
			w.WriteHeader(http.StatusUnauthorized)
			return
		case err != nil:
			fmt.Println("Query error: %v\n", err)
			// return 500 status
			w.WriteHeader(http.StatusInternalServerError)
			return
		default:
			res.Logged_In = true
			res.Message = "success"
			json.NewEncoder(w).Encode(res)
	}
}

func createToken(username string, db *sql.DB) (string, int) {
	var user_id int
	if err := db.QueryRow("select user_id from credentials where username = $1", username).Scan(&user_id); err != nil {
		panic(err)
	}
	type Claims struct {
		Username string `json:"username"`
		User_Id int `json:"user_id"`
		jwt.StandardClaims
	}
	claims := Claims{
		username,
		user_id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(365 * 24 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(getSecretKey())
	if err != nil {
		panic(err)
	}
	return ss, user_id
}

func getSecretKey() ([]byte) {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}
	key := []byte(os.Getenv("SECRET_KEY"))
	return key
}
