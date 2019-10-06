package attendees

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Child struct {
	gorm.Model
	FirstName string		`gorm:"size:100"`
	DateOfBirth time.Time
	Guardian Guardian
}