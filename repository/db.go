package repository

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
	"time"
)

var db *gorm.DB
var err error
var deleteOnExit bool

func InitDB() {
	time.Sleep(15 * time.Second) //give MySql time to start up
	db, err = gorm.Open("mysql",  os.Getenv("GORM_CONNSTR"))
	if err != nil {
		fmt.Printf("failed to connect database: %v", err)
		panic("failed to connect database")
	}
	db.LogMode(true)

	// Migrate the schema
	db.AutoMigrate(&Event{})

	// Load 6 months of past events
	row := Event{}
	if db.First(&row).RowsAffected < 1 {
		now := time.Now()
		now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		for i := 1; i < 183; i++ {
			db.Create(&Event{EventDate: now.AddDate(0, 0, -(i)), Active: false, RegistrationOpen: false})
		}
		// load a future event1 that is active
		db.Create(&Event{EventDate: time.Date(2021, 5, 20, 0, 0, 0, 0, time.UTC), Active: true, RegistrationOpen: false})
		deleteOnExit = true
	}

}

func Cleanup() {
	// Delete - delete everything
	if deleteOnExit {
		db.Unscoped().Delete(Event{})
		if err := db.Close(); err != nil {
			fmt.Printf("Error on database close: %v", err)
		}
		fmt.Println("Completed Deleting Events")
	}
}