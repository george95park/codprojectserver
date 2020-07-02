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
		res := models.Error {
			Status: "400",
			Text: "Status Bad Request",
			Message: "Bad request",
		}
		json.NewEncoder(w).Encode(res)
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
		res := models.Error {
			Status: "500",
			Text: "Status Internal Server Error",
			Message: "Username already taken",
		}
		json.NewEncoder(w).Encode(res)
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
		HttpOnly: true,
		//Secure: true,
	})
	user := models.User{
		Username: creds.Username,
		User_Id: user_id,
		Logged_In: true,
	}
	json.NewEncoder(w).Encode(user)
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
		res := models.Error {
			Status: "400",
			Text: "Status Bad Request",
			Message: "Bad request",
		}
		json.NewEncoder(w).Encode(res)
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
			res := models.Error {
				Status: "401",
				Text: "Status Unauthorized",
				Message: "No user with this username",
			}
			json.NewEncoder(w).Encode(res)
			return
		case err != nil:
			fmt.Println("Query error: %v\n", err)
			// return 500 status
			w.WriteHeader(http.StatusInternalServerError)
			res := models.Error {
				Status: "500",
				Text: "Status Internal Server Error",
				Message: "Internal server error",
			}
			json.NewEncoder(w).Encode(res)
			return
		default:
			if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password)); err != nil {
				fmt.Println("Access Denied: Wrong Password.")
				// return 401 status
				w.WriteHeader(http.StatusUnauthorized)
				res := models.Error {
					Status: "401",
					Text: "Status Unauthorized",
					Message: "Wrong password",
				}
				json.NewEncoder(w).Encode(res)
				return
			} else {
				fmt.Println("Access granted.")
				token, user_id := createToken(creds.Username, db)
				_, err := db.Exec("update credentials set token=$1 where user_id=$2", token, user_id)
				if err != nil {
					fmt.Println(err)

					// return 500 status
					w.WriteHeader(http.StatusInternalServerError)
					res := models.Error {
						Status: "500",
						Text: "Status Internal Server Error",
						Message: "Internal server error",
					}
					json.NewEncoder(w).Encode(res)
					return
				}
				http.SetCookie(w, &http.Cookie {
					Name: "token",
					Value: token,
					Expires: time.Now().Add(365 * 24 * time.Hour),
					HttpOnly: true,
					//Secure: true,
				})
				user := models.User {
					Username: creds.Username,
					User_Id: user_id,
					Logged_In: true,
				}
				json.NewEncoder(w).Encode(user)
			}
	}
}
func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	http.SetCookie(w, &http.Cookie {
		Name: "token",
		Value: "",
		Expires: time.Now(),
		HttpOnly: true,
		//Secure: true,
	})
	fmt.Println("Logged out")
}

func GetSessionTokenUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
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
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			user := models.User{
				Username: "",
				User_Id: -1,
				Logged_In: false,
			}
			json.NewEncoder(w).Encode(user)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		res := models.Error {
			Status: "400",
			Text: "Status Bad Request",
			Message: "Bad request",
		}
		json.NewEncoder(w).Encode(res)
		return
	}
	tokenStr := c.Value
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token * jwt.Token) (interface{}, error) {
		return getSecretKey(), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			res := models.Error {
				Status: "401",
				Text: "Status Unauthorized",
				Message: "Invalid",
			}
			json.NewEncoder(w).Encode(res)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		res := models.Error {
			Status: "400",
			Text: "Status Bad Request",
			Message: "Bad request",
		}
		json.NewEncoder(w).Encode(res)
		return
	}
	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		res := models.Error {
			Status: "401",
			Text: "Status Unauthorized",
			Message: "Invalid",
		}
		json.NewEncoder(w).Encode(res)
		return
	}
	user := models.User{
		Username: claims.Username,
		User_Id: claims.User_Id,
		Logged_In: true,
	}
	json.NewEncoder(w).Encode(user)
}

func createToken(username string, db *sql.DB) (string, int) {
	var user_id int
	if err := db.QueryRow("select user_id from credentials where username = $1", username).Scan(&user_id); err != nil {
		panic(err)
	}
	claims := models.Claims{
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
