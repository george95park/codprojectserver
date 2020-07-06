package config

import (
	"fmt"
	"codproject/server/lib"
)

func CreateTables() {
    db := lib.ConnectDB()
    defer db.Close()
	const gunsSchema =
	`CREATE TABLE IF NOT EXISTS guns (
		gun_id SERIAL PRIMARY KEY,
		type TEXT,
		name TEXT
	)`
	const attachmentsSchema =
	`CREATE TABLE IF NOT EXISTS attachments (
		attachment_id SERIAL PRIMARY KEY,
		gun_id INT,
		name TEXT,
		subattachments TEXT[],
		FOREIGN KEY (gun_id) REFERENCES guns(gun_id)
	)`
	const credentialsSchema =
	`CREATE TABLE IF NOT EXISTS credentials (
		user_id SERIAL PRIMARY KEY,
		password TEXT,
		username TEXT UNIQUE,
		token TEXT
	)`
	const loadoutsSchema =
	`CREATE TABLE IF NOT EXISTS loadouts (
		loadout_id SERIAL PRIMARY KEY,
		user_id INT,
		type TEXT,
		gun TEXT,
		attachments TEXT[],
		subattachments TEXT[],
		description TEXT,
		FOREIGN KEY (user_id) REFERENCES credentials(user_id)
	)`

	queryStrings := []string{gunsSchema, attachmentsSchema, credentialsSchema, loadoutsSchema}
	for i := 0; i < len(queryStrings); i++ {
		fmt.Println("Executing query string: ", i)
		if _,err := db.Exec(queryStrings[i]); err != nil {
			panic(err)
		}
		fmt.Println("Done.")
	}
}
