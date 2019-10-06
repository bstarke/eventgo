package main

import (
	"encoding/json"
	"eventgo/event"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var db *gorm.DB
var err error

func main() {
	db, err = gorm.Open("sqlite3", "events.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer cleanup()
	db.LogMode(true)

	// Migrate the schema
	db.AutoMigrate(&event.Event{})

	// Load a year of past events
	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	for i := 0; i < 1; i++ {
		db.Create(&event.Event{EventDate: now.AddDate(0, 0, -(i)), Active: false, RegistrationOpen: false})
	}
	// load a future event1 that is active
	db.Create(&event.Event{EventDate: time.Date(2021, 5, 20, 0, 0, 0, 0, time.UTC), Active: true, RegistrationOpen: false})

	handleRequests()

}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/events", allEvents)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func homePage(w http.ResponseWriter, r *http.Request){
	_, _ = fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func cleanup() {
	// Delete - delete everything
	db.Unscoped().Delete(event.Event{})
	_ = db.Close()
}

func allEvents(w http.ResponseWriter, r *http.Request) {
	events := []event.Event{}
	db.Find(&events)
	_ = json.NewEncoder(w).Encode(&events)
}

func createNewEvent(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// return the string response containing the request body
	reqBody, _ := ioutil.ReadAll(r.Body)
	var newEvent event.Event
	_ = json.Unmarshal(reqBody, &newEvent)
	db.Create(&newEvent)
	_ = json.NewEncoder(w).Encode(newEvent)
}

func test() {
	// Read
	var event1 event.Event
	db.First(&event1, 220) // find event1 with id ?
	fmt.Printf("Event: %v\n", event1)
	event1 = event.Event{}
	db.First(&event1, "active = ?", true) // find event1 with code l1212
	fmt.Printf("Event: %v\n", event1)

	// Update - update event1's registration to true
	db.Model(&event1).Update("RegistrationOpen", true)
	fmt.Printf("Event: %v\n", event1)
}