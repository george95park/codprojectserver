package config

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

func InitDB() {
    db := ConnectDB()
    defer db.Close()
	queryStrings := []string{
		"create table guns (gun_id serial primary key, type text, name text)",
		"create table attachments (attachment_id serial primary key, gun_id int, name text, subattachments text[], foreign key (gun_id) references guns(gun_id))",
		"create table credentials (user_id serial primary key, password text, username text, token text)",
		"create table loadouts (loadout_id serial primary key, user_id int, type text, gun text, attachments text[], subattachments text[], description text, foreign key (user_id) references credentials(user_id))",
	}
	for i := 0; i < len(queryStrings); i++ {
		fmt.Println("Executing query string: ", queryStrings[i])
		if _,err := db.Exec(queryStrings[i]); err != nil {
			panic(err)
		}
		fmt.Println("Done.")
	}
}
