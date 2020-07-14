package middleware

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/lib/pq"
	"codproject/server/models"
	"codproject/server/lib"
)

func GetGuns(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "http://codloadish.com")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	db := lib.ConnectDB()
	defer db.Close()

	g := &models.Gun{}
	if err := json.NewDecoder(r.Body).Decode(g); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rows,err := db.Query("select * from guns where type=$1", g.Type)
	if err != nil {
		panic(err)
	}
	allGuns := []models.Gun{}
	for rows.Next() {
		gun := models.Gun{}
		err = rows.Scan(&gun.Gun_Id, &gun.Type, &gun.Name)
		if err != nil {
			panic(err)
		}
		allGuns = append(allGuns, gun)
	}
	fmt.Println("Received guns.")

	// return response
	json.NewEncoder(w).Encode(allGuns)
}

func GetAttachments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "http://codloadish.com")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	db := lib.ConnectDB()
	defer db.Close()

	a := &models.Attachment{}
	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rows,err := db.Query("select * from attachments where gun_id=$1", a.Gun_Id)
	allAttachments := []models.Attachment{}

	for rows.Next() {
		att := models.Attachment{}
		err = rows.Scan(&att.Attachment_Id, &att.Gun_Id, &att.Name, pq.Array(&att.SubAttachments))
		if err != nil {
			panic(err)
		}
		allAttachments = append(allAttachments, att)
	}
	fmt.Println("Received Attachments.")
	json.NewEncoder(w).Encode(allAttachments)
}
