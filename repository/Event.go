package repository

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Event struct {
	gorm.Model
	EventDate time.Time
	Active bool
	RegistrationOpen bool
}

func (e Event) FindAll() (events []Event){
	db.Find(&events)
	return
}

func (e Event) FindById(id int64) (event Event) {
	db.Find(&event, id)
	return
}

func (e Event) Create(newEvent *Event) (event Event){
	db.Create(&newEvent)
	event = *newEvent
	return
}