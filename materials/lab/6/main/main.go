package main

import (
	"log"
	"net/http"
	"scrape/scrape"

	"github.com/gorilla/mux"
)

// LOG_LEVEL defined in logging.go

func main() {

	log.Println("starting API server")
	//create a new router
	router := mux.NewRouter()
	log.Println("creating routes")
	//specify endpoints
	router.HandleFunc("/", scrape.MainPage).Methods("GET")

	router.HandleFunc("/api-status", scrape.APISTATUS).Methods("GET")

	router.HandleFunc("/indexer", scrape.IndexFiles).Methods("GET")
	router.HandleFunc("/search", scrape.FindFile).Methods("GET")
	router.HandleFunc("/addsearch/{regex}", scrape.AddSearch).Methods("GET")
	router.HandleFunc("/clear", scrape.Clear).Methods("GET")
	router.HandleFunc("/reset", scrape.Reset).Methods("GET")

	http.Handle("/", router)

	//start and listen to requests
	http.ListenAndServe(":8080", router)

}
