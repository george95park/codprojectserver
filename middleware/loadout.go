package middleware

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/lib/pq"
	"github.com/gorilla/mux"
	"codproject/server/models"
	"codproject/server/config"
)

// creates loadout in database
func CreateLoadout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	db := config.ConnectDB()
	defer db.Close()
	load := &models.Loadout{}
	// Unmarshall request body
	if err := json.NewDecoder(r.Body).Decode(load); err != nil {
		panic(err)
	}

	// insert new loadout to database
	if _, err := db.Query("insert into loadouts (user_id, type, gun, attachments, subattachments, description) values ($1, $2, $3, $4, $5, $6)",
		load.User_Id,
		load.Type,
		load.Gun,
		pq.Array(load.Attachments),
		pq.Array(load.SubAttachments),
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
	db := config.ConnectDB()
	defer db.Close()
	currUserId := mux.Vars(r)["id"]

	// Get the rows according to user
	rows,err := db.Query("select * from loadouts where user_id=$1",currUserId)
	if err != nil {
		panic(err)
	}

	// For each row, create new loadout object and append to the result
	userLoadouts := []models.Loadout{}
	for rows.Next() {
		l := models.Loadout{}
		err = rows.Scan(&l.Loadout_Id, &l.User_Id, &l.Type, &l.Gun, pq.Array(&l.Attachments), pq.Array(&l.SubAttachments), &l.Description)
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
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// open database
	db := config.ConnectDB()
	defer db.Close()
	currLoadoutId := mux.Vars(r)["id"]

	res,err := db.Exec("delete from loadouts where loadout_id=$1", currLoadoutId)
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
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	// open database
	db := config.ConnectDB()
	defer db.Close()
	load := &models.Loadout{}
	if err := json.NewDecoder(r.Body).Decode(load); err != nil {
		fmt.Println(err)
		// return a 400 status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res, err := db.Exec("update loadouts set type=$2, gun=$3, attachments=$4, subattachments=$5, description=$6 where loadout_id=$1",
		load.Loadout_Id,
		load.Type,
		load.Gun,
		pq.Array(load.Attachments),
		pq.Array(load.SubAttachments),
		load.Description)
	if err != nil {
		panic(err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Println("Total rows affected: %v", rowsAffected)
	json.NewEncoder(w).Encode(load)
}
