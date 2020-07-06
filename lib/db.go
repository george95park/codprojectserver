package lib

import (
	"os"
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
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
