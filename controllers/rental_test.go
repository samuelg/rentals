package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/samuelg/rentals/config"
	log "github.com/samuelg/rentals/logging"
	"github.com/stretchr/testify/suite"
)

// Test suite for the Rental controller
type RentalControllerTestSuite struct {
	suite.Suite
	config *config.Config
	router *gin.Engine
}

func setupRouter() *gin.Engine {
	// don't log API calls
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	rentalGroup := router.Group("rentals")
	{
		rentals := new(RentalController)
		rentalGroup.GET("/", rentals.List)
		rentalGroup.GET("/:rental_id", rentals.Get)
	}

	return router
}

func (suite *RentalControllerTestSuite) SetupSuite() {
	config.Init("test")
	log.Init("FATAL", config.GetConfig().AppVersion)
	suite.config = config.GetConfig()
	suite.router = setupRouter()
}

// used to unmarshal json responses
type listResponse struct {
	Data []string `json:"data"`
}

// used to unmarshal json responses
type getResponse struct {
	Id int32 `json:"id"`
}

// GET /rentals tests

func (suite *RentalControllerTestSuite) TestListRentalsSuccess() {
	req, _ := http.NewRequest("GET", "/rentals/", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	if suite.Equal(http.StatusOK, w.Code) {
		var response listResponse
		if suite.Nil(json.Unmarshal(w.Body.Bytes(), &response), "Should be able to unmarshal response") {
			suite.Equal(0, len(response.Data), "Should return an empty array")
		}
	}
}

func (suite *RentalControllerTestSuite) TestListRentalsSuccessAllFilters() {
	req, _ := http.NewRequest("GET", "/rentals/?price_min=100&price_max=200000&limit=10&offset=0&sort=price&ids=1,2&near=61.0,-149.0", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	if suite.Equal(http.StatusOK, w.Code) {
		var response listResponse
		if suite.Nil(json.Unmarshal(w.Body.Bytes(), &response), "Should be able to unmarshal response") {
			suite.Equal(0, len(response.Data), "Should return an empty array")
		}
	}
}

func (suite *RentalControllerTestSuite) TestListRentalsInvalidPriceMin() {
	req, _ := http.NewRequest("GET", "/rentals/?price_min=a", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *RentalControllerTestSuite) TestListRentalsInvalidPriceMax() {
	req, _ := http.NewRequest("GET", "/rentals/?price_max=a", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *RentalControllerTestSuite) TestListRentalsInvalidLimit() {
	req, _ := http.NewRequest("GET", "/rentals/?limit=50.5", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *RentalControllerTestSuite) TestListRentalsLimitTooLarge() {
	req, _ := http.NewRequest("GET", "/rentals/?limit=1000", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *RentalControllerTestSuite) TestListRentalsInvalidOffset() {
	req, _ := http.NewRequest("GET", "/rentals/?offset=50.5", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *RentalControllerTestSuite) TestListRentalsInvalidSort() {
	req, _ := http.NewRequest("GET", "/rentals/?sort=notafield", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *RentalControllerTestSuite) TestListRentalsInvalidIds() {
	req, _ := http.NewRequest("GET", "/rentals/?ids=1,a", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *RentalControllerTestSuite) TestListRentalsInvalidNear() {
	req, _ := http.NewRequest("GET", "/rentals/?near=1,a", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *RentalControllerTestSuite) TestListRentalsNearTooFewValues() {
	req, _ := http.NewRequest("GET", "/rentals/?near=1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *RentalControllerTestSuite) TestListRentalsNearTooManyValues() {
	req, _ := http.NewRequest("GET", "/rentals/?near=1,2,3", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *RentalControllerTestSuite) TestListRentalsNearInvalidLatLong() {
	req, _ := http.NewRequest("GET", "/rentals/?near=100,200", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
}

// GET /rentals/:rental_id tests

func (suite *RentalControllerTestSuite) TestGetRentalSuccess() {
	req, _ := http.NewRequest("GET", "/rentals/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	if suite.Equal(http.StatusOK, w.Code) {
		var response getResponse
		if suite.Nil(json.Unmarshal(w.Body.Bytes(), &response), "Should be able to unmarshal response") {
			suite.Equal(int32(1), response.Id)
		}
	}
}

func (suite *RentalControllerTestSuite) TestGetRentalInvalidId() {
	req, _ := http.NewRequest("GET", "/rentals/invalid", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
}

// test for invalid route handling

func (suite *RentalControllerTestSuite) TestRouteNotFound() {
	req, _ := http.NewRequest("GET", "/invalid/route", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusNotFound, w.Code)
}

func TestRentalControllerTestSuite(t *testing.T) {
	suite.Run(t, new(RentalControllerTestSuite))
}
