package main

import (
	log "github.com/sirupsen/logrus"
	"gorilla/mux"
	"net/http"
)

func main() {
	route := mux.NewRouter()

	// base path
	prefix := route.PathPrefix("/api").Subrouter()

	// routes
	prefix.HandleFunc("/createProfile", createProfile).Methods("POST")
	prefix.HandleFunc("/getAllStudents", getAllStudents).Methods("GET")
	prefix.HandleFunc("/getStudentProfile", getStudentProfile).Methods("POST")
	prefix.HandleFunc("/updateProfile", updateProfile).Methods("PUT")
	prefix.HandleFunc("/deleteProfile/{id}", deleteProfile).Methods("DELETE")
	
	// run server
	log.Fatal(http.ListenAndServe(":9000", prefix))
}