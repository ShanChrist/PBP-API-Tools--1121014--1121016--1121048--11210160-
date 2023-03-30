package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ExplorasiTools/controllers"
	"github.com/go-co-op/gocron"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	c := gocron.NewScheduler(time.UTC)

	router.HandleFunc("/users", controllers.GetAllUsers).Methods("GET")
	router.HandleFunc("/forgot_password", controllers.ForgotPassword).Methods("GET")
	router.HandleFunc("/cron", controllers.Cron).Methods("GET")
	c.StartAsync()

	router.HandleFunc("/login", controllers.Login).Methods("POST")
	router.HandleFunc("/logout", controllers.Logout).Methods("GET")

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":7070", router))
}
