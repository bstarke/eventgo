package main

import (
	"encoding/json"
	"eventgo/repository"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		repository.Cleanup()
		fmt.Println("Exiting")
		os.Exit(0)
	}()
	repository.InitDB()
	handleRequests()
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.Use(commonMiddleware)
	myRouter.HandleFunc("/", homePage).Methods(http.MethodGet)
	myRouter.HandleFunc("/events", allEvents).Methods(http.MethodGet)
	myRouter.HandleFunc("/events", createEvent).Methods(http.MethodPost)
	myRouter.HandleFunc("/events/active", eventByActive).Methods(http.MethodGet)
	myRouter.HandleFunc("/events/{id}", eventById).Methods(http.MethodGet)
	myRouter.HandleFunc("/events/{id}", updateEvent).Methods(http.MethodPut, http.MethodPatch)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func homePage(w http.ResponseWriter, r *http.Request) {
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
	event := repository.Event{}.FindById(uint(id))
	_ = json.NewEncoder(w).Encode(&event)
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var event repository.Event
	_ = json.Unmarshal(reqBody, &event)
	repository.Event{}.Create(&event)
	_ = json.NewEncoder(w).Encode(&event)
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var event map[string]interface{}
	_ = json.Unmarshal(reqBody, &event)
	vars := mux.Vars(r)
	uId, _ := strconv.Atoi(vars["id"])
	nEvent, _ := repository.Event{}.PatchUpdate(uint(uId), event)
	_ = json.NewEncoder(w).Encode(&nEvent)
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}