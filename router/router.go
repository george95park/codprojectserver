package router

import (
	"codproject/server/middleware"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", middleware.HomePage)
	router.HandleFunc("/login", middleware.Login)
	router.HandleFunc("/signup", middleware.Signup)
	return router
}
