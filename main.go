package main

import (
	"encoding/json"
	"eventgo/repository"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)


func main() {
	repository.InitDB()
	handleRequests()
	defer repository.Cleanup()
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage).Methods(http.MethodGet)
	myRouter.HandleFunc("/events", allEvents).Methods(http.MethodGet)
	myRouter.HandleFunc("/events", createEvent).Methods(http.MethodPost)
	myRouter.HandleFunc("/events/active", eventByActive).Methods(http.MethodGet)
	myRouter.HandleFunc("/events/{id}", eventById).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func homePage(w http.ResponseWriter, r *http.Request){
	_, _ = fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func allEvents(w http.ResponseWriter, r *http.Request) {
	events := repository.Event{}.FindAll()
	_ = json.NewEncoder(w).Encode(&events)
}

func eventByActive(w http.ResponseWriter, r *http.Request) {
	var event = repository.Event{}.FindActive()
	_ = json.NewEncoder(w).Encode(&event)
}

func eventById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	id, _ := strconv.Atoi(key)
	event := repository.Event{}.FindById(int64(id))
	_ = json.NewEncoder(w).Encode(&event)
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var event repository.Event
	_ = json.Unmarshal(reqBody, &event)

	_ = json.NewEncoder(w).Encode(&event)
}