package lib

import (
	"os"
	"time"
	"database/sql"
	"github.com/joho/godotenv"
	"codproject/server/models"
	"github.com/dgrijalva/jwt-go"
)

func CreateToken(username string, db *sql.DB) (string, int) {
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
	ss, err := token.SignedString(GetSecretKey())
	if err != nil {
		panic(err)
	}
	return ss, user_id
}

func GetSecretKey() ([]byte) {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}
	key := []byte(os.Getenv("SECRET_KEY"))
	return key
}
