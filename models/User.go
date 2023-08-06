package models

type User struct {
	ID        int32    `gorm:"primary_key;autoincrement;column:id"`
	FirstName string   `gorm:"column:first_name"`
	LastName  string   `gorm:"column:last_name"`
	Rentals   []Rental `gorm:"foreignKey:UserId"`
}
