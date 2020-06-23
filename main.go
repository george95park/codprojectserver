package main

import (
	"fmt"
	"log"
	"net/http"
	"codproject/server/router"
	"codproject/server/config"
)

func main() {
	config.CreateTables()
	r := router.Router()
	fmt.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
