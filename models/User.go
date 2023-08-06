package models

// User model
type User struct {
	// uses serial integer column in the database which will use 32 bits at most
	ID        uint32   `gorm:"primary_key;autoincrement;column:id"`
	FirstName string   `gorm:"column:first_name"`
	LastName  string   `gorm:"column:last_name"`
	Rentals   []Rental `gorm:"foreignKey:UserId"`
}
