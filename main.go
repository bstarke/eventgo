package main

import (
	"eventgo/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var version Version
var GitHash string
var BuildTime string
var GoVer string

type Version struct {
	GitCommit  string
	ApiVersion string
	GoVersion  string
	BuildDate  string
}

func main() {
	version = Version{GitCommit: GitHash, ApiVersion: "1.0.0", BuildDate: BuildTime, GoVersion: GoVer}
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
	r := gin.Default()
	//r.Use(commonMiddleware)
	r.GET("/", homePage)
	r.GET("/events", allEvents)
	r.POST("/events", createEvent)
	r.GET("/events/active", eventByActive)
	r.GET("/events/:id", eventById)
	r.PUT("/events/:id", updateEvent)
	r.PATCH("/events/:id", updateEvent)
	log.Fatal(r.Run(":8080"))
}

func homePage(c *gin.Context) {
	fmt.Println("Endpoint Hit: homePage")
	c.JSON(http.StatusOK, &version)
}

func allEvents(c *gin.Context) {
	events := repository.Event{}.FindAll()
	c.JSON(http.StatusOK, &events)
}

func eventByActive(c *gin.Context) {
	var event = repository.Event{}.FindActive()
	c.JSON(http.StatusOK, &event)
}

func eventById(c *gin.Context) {
	key := c.Param("id")
	id, _ := strconv.Atoi(key)
	event := repository.Event{}.FindById(uint(id))
	c.JSON(http.StatusOK, &event)
}

func createEvent(c *gin.Context) {
	var event repository.Event
	if c.ShouldBind(&event) == nil {
		repository.Event{}.Create(&event)
	}
	c.JSON(http.StatusOK, &event)
}

func updateEvent(c *gin.Context) {
	var event map[string]interface{}
	err := c.BindJSON(&event)
	if err != nil {
		log.Printf("Error Binding JSON.. haha %s", err.Error())
		return
	}
	uId, _ := strconv.Atoi(c.Param("id"))
	nEvent, _ := repository.Event{}.PatchUpdate(uint(uId), event)
	c.JSON(http.StatusOK, &nEvent)
}
