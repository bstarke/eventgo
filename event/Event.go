package event

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