package middleware

import (
	"fmt"
	"net/http"
	"encoding/json"
	"codproject/server/models"
	"codproject/server/lib"
)


func GetAllUsers(w http.ResponseWriter, r * http.Request) {
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

	rows,err := db.Query("select user_id, username from credentials")
	if err != nil {
		panic(err)
	}
	res := []models.User{}
	for rows.Next() {
		r := models.User{}
		err = rows.Scan(&r.User_Id, &r.Username)
		if err != nil {
			panic(err)
		}
		res = append(res, r)
	}
	fmt.Println("Sending all users")
	json.NewEncoder(w).Encode(res)

}
