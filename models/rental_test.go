package models

import (
	"encoding/json"
	"testing"

	"github.com/samuelg/rentals/config"
	"github.com/samuelg/rentals/db"
	log "github.com/samuelg/rentals/logging"
	"github.com/stretchr/testify/suite"
)

// Test suite for the Rental controller
type RentalModelTestSuite struct {
	suite.Suite
	config *config.Config
}

func (suite *RentalModelTestSuite) SetupSuite() {
	config.Init("test")
	log.Init("FATAL", config.GetConfig().AppVersion)
	db.Init()
	suite.config = config.GetConfig()
}

func (suite *RentalModelTestSuite) TestMarshallJson() {
	rental := Rental{
		ID:          1,
		Name:        "My rental",
		Description: "My awesome rental",
		User: User{
			ID:        1,
			FirstName: "Bob",
			LastName:  "Smith",
		},
		Type:            "RV",
		Sleeps:          3,
		Price:           1000,
		City:            "Huntsville",
		State:           "AL",
		Zip:             "35758",
		Country:         "US",
		VehicleMake:     "VW",
		VehicleModel:    "Bus",
		VehicleYear:     1970,
		VehicleLength:   15.50,
		Lat:             36.1,
		Lng:             -86.4,
		PrimaryImageUrl: "http://images.com/1.png",
	}

	bytes, err := json.Marshal(rental)
	suite.Nil(err, "Should be able to marshal")
	suite.Equal(`{"id":1,"name":"My rental",`+
		`"description":"My awesome rental",`+
		`"type":"RV","make":"VW","model":"Bus",`+
		`"year":1970,"length":15.5,"sleeps":3,`+
		`"primary_image_url":"http://images.com/1.png","price":{"day":1000},`+
		`"location":{"city":"Huntsville","state":"AL","zip":"35758",`+
		`"country":"US","lat":36.1,"lng":-86.4},"user":{"id":1,`+
		`"first_name":"Bob","last_name":"Smith"}}`, string(bytes))
}

func TestRentalModelTestSuite(t *testing.T) {
	suite.Run(t, new(RentalModelTestSuite))
}
