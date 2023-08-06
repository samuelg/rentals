package models

import (
	"encoding/json"
	"time"
)

type Rental struct {
	ID              int32 `gorm:"primary_key;autoincrement;column:id"`
	UserId          int32 `gorm:"column:user_id"`
	User            User
	Name            string    `gorm:"column:name"`
	Type            string    `gorm:"column:type"`
	Description     string    `gorm:"column:description"`
	Sleeps          int32     `gorm:"column:sleeps"`
	Price           int64     `gorm:"column:price_per_day"`
	City            string    `gorm:"column:home_city"`
	State           string    `gorm:"column:home_state"`
	Zip             string    `gorm:"column:home_zip"`
	Country         string    `gorm:"column:home_country"`
	VehicleMake     string    `gorm:"column:vehicle_make"`
	VehicleModel    string    `gorm:"column:vehicle_model"`
	VehicleYear     int32     `gorm:"column:vehicle_year"`
	VehicleLength   float32   `gorm:"column:vehicle_length;precision:4;scale:2"`
	Created         time.Time `gorm:"column:created"`
	Updated         time.Time `gorm:"column:updated"`
	Lat             float32   `gorm:"column:lat"`
	Lng             float32   `gorm:"column:lng"`
	PrimaryImageUrl string    `gorm:"column:primary_image_url"`
}

// Custom JSON format for the response
func (rental Rental) MarshalJSON() ([]byte, error) {
	type Price struct {
		Day int64 `json:"day"`
	}
	type Location struct {
		City    string  `json:"city"`
		State   string  `json:"state"`
		Zip     string  `json:"zip"`
		Country string  `json:"country"`
		Lat     float32 `json:"lat"`
		Lng     float32 `json:"lng"`
	}
	type User struct {
		Id        int32  `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	return json.Marshal(&struct {
		ID              int32    `json:"id"`
		Name            string   `json:"name"`
		Description     string   `json:"description"`
		Type            string   `json:"type"`
		Make            string   `json:"make"`
		Model           string   `json:"model"`
		Year            int32    `json:"year"`
		Length          float32  `json:"length"`
		Sleeps          int32    `json:"sleeps"`
		PrimaryImageUrl string   `json:"primary_image_url"`
		Price           Price    `json:"price"`
		Location        Location `json:"location"`
		User            User     `json:"user"`
	}{
		ID:              rental.ID,
		Name:            rental.Name,
		Description:     rental.Description,
		Type:            rental.Type,
		Make:            rental.VehicleMake,
		Model:           rental.VehicleModel,
		Year:            rental.VehicleYear,
		Length:          rental.VehicleLength,
		Sleeps:          rental.Sleeps,
		PrimaryImageUrl: rental.PrimaryImageUrl,
		Price:           Price{Day: rental.Price},
		Location: Location{
			City:    rental.City,
			State:   rental.State,
			Zip:     rental.Zip,
			Country: rental.Country,
			Lat:     rental.Lat,
			Lng:     rental.Lng,
		},
		User: User{
			Id:        rental.User.ID,
			FirstName: rental.User.FirstName,
			LastName:  rental.User.LastName,
		},
	})
}
