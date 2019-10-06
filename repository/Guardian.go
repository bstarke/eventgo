package repository
import (
	"github.com/jinzhu/gorm"
)

type Guardian struct {
	gorm.Model
	FirstName string	`gorm:"size:100"`
	LastName string		`gorm:"size:100"`
	Address string		`gorm:"size:255"`
	City string			`gorm:"size:100"`
	ZipCode int8
	Children []Child
}