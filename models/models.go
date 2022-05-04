package models

type Users struct {
	ID        uint   `gorm:"primary_key"`
	Username  string `gorm:"not null;unique"`
	Password  string `gorm:"not null,size:255"`
	Email     string `gorm:"not null;unique"`
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`
}
