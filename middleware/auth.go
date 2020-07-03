package middleware

import (
	"net/http"
	"encoding/json"
	"codproject/server/models"
	"codproject/server/config"
    "codproject/server/lib"
	"github.com/dgrijalva/jwt-go"
)

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
		return lib.GetSecretKey(), nil
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
