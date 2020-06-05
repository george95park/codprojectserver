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
	router.HandleFunc("/createloadout", middleware.CreateLoadout)
	router.HandleFunc("/getloadouts", middleware.GetLoadouts)
	router.HandleFunc("/deleteloadout", middleware.DeleteLoadout)
	router.HandleFunc("/updateloadout", middleware.UpdateLoadout)
	return router
}
