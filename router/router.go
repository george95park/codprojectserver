package router

import (
	"codproject/server/middleware"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/login", middleware.Login)
	router.HandleFunc("/signup", middleware.Signup)
	router.HandleFunc("/createloadout", middleware.CreateLoadout)
	router.HandleFunc("/getloadouts/{id}", middleware.GetLoadouts)
	router.HandleFunc("/deleteloadout/{id}", middleware.DeleteLoadout)
	router.HandleFunc("/updateloadout", middleware.UpdateLoadout)
	router.HandleFunc("/getsessiontokenuser", middleware.GetSessionTokenUser)
	router.HandleFunc("/getguns", middleware.GetGuns)
	router.HandleFunc("/getattachments", middleware.GetAttachments)
	router.HandleFunc("/getallusers", middleware.GetAllUsers)
	return router
}
