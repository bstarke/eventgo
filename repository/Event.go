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

func (e Event) FindById(id uint) (event Event) {
	db.Find(&event, id)
	return
}

func (e Event) FindActive() (event Event) {
	db.Where("active = ?", true).First(&event)
	return
}

func (e Event) Create(newEvent *Event) (event Event){
	db.Create(&newEvent)
	event = *newEvent
	return
}

func (e Event) Update(uEvent *Event) error {
	db.Model(&uEvent).Updates(&uEvent)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (e Event) PatchUpdate(uid uint, uEvent map[string]interface{}) (event Event, err error) {
	db.First(&event, uid)
	db.Model(&event).Updates(uEvent)
	if db.Error != nil {
		return Event{}, db.Error
	}
	return event, nil
}