package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"database/sql"
	"encoding/json"
	"github.com/lib/pq"
	"github.com/gorilla/mux"
	"codproject/server/models"
	"codproject/server/lib"
)

// creates loadout in database
func CreateLoadout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "http://codloadish.com")
	//w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	db := lib.ConnectDB()
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
	w.Header().Set("Access-Control-Allow-Origin", "http://codloadish.com")
	//w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	db := lib.ConnectDB()
	defer db.Close()
	currUserId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err)
	}
	// get user loadouts
	userLoadouts := getUserLoadouts(db, currUserId)


	// return response
	json.NewEncoder(w).Encode(userLoadouts)
}

func DeleteLoadout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "http://codloadish.com")
	//w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// open database
	db := lib.ConnectDB()
	defer db.Close()
	currLoadoutId := mux.Vars(r)["id"]
	fmt.Println(currLoadoutId)
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
	// return response
	json.NewEncoder(w).Encode(currLoadoutId)
}

func UpdateLoadout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "http://codloadish.com")
	//w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	// open database
	db := lib.ConnectDB()
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
	// return response
	json.NewEncoder(w).Encode(load)
}

func getUserLoadouts(db *sql.DB, userid int) []models.Loadout {
	// Get the rows according to user
	rows,err := db.Query("select * from loadouts where user_id=$1",userid)
	if err != nil {
		panic(err)
	}

	// For each row, create new loadout object and append to the result
	res := []models.Loadout{}
	for rows.Next() {
		l := models.Loadout{}
		err = rows.Scan(&l.Loadout_Id, &l.User_Id, &l.Type, &l.Gun, pq.Array(&l.Attachments), pq.Array(&l.SubAttachments), &l.Description)
		if err != nil {
			panic(err)
		}
		res = append(res, l)
	}
	fmt.Println("Received loadouts.")
	return res
}
