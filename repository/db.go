package repository

import (
	"github.com/jinzhu/gorm"
	"time"
)

var db *gorm.DB
var err error

func InitDB() {
	db, err = gorm.Open("sqlite3", "events.db")
	if err != nil {
		panic("failed to connect database")
	}
	db.LogMode(true)

	// Migrate the schema
	db.AutoMigrate(&Event{})

	// Load a year of past events
	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	for i := 0; i < 1; i++ {
		db.Create(&Event{EventDate: now.AddDate(0, 0, -(i)), Active: false, RegistrationOpen: false})
	}
	// load a future event1 that is active
	db.Create(&Event{EventDate: time.Date(2021, 5, 20, 0, 0, 0, 0, time.UTC), Active: true, RegistrationOpen: false})

}

func Cleanup() {
	// Delete - delete everything
	db.Unscoped().Delete(Event{})
	_ = db.Close()
}