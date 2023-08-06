package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/samuelg/rentals/config"
	"github.com/samuelg/rentals/db"
	log "github.com/samuelg/rentals/logging"
	"github.com/samuelg/rentals/models"
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
	db.Init()
	suite.config = config.GetConfig()
	suite.router = setupRouter()
}

// GET /rentals tests
type testListResponse struct {
	Pagigation *PaginationResponse `json:"pagination"`
	// after a Rental model is marshaled
	Data []models.RentalResponse `json:"data"`
}

func (suite *RentalControllerTestSuite) TestListRentalsSuccess() {
	req, _ := http.NewRequest("GET", "/rentals/?limit=1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	if suite.Equal(http.StatusOK, w.Code) {
		var response testListResponse
		if suite.Nil(json.Unmarshal(w.Body.Bytes(), &response), "Should be able to unmarshal response") {
			suite.Equal(1, len(response.Data), "Should return a single result")
			suite.Equal(uint32(30), response.Pagigation.Count)
			suite.Equal(uint8(1), response.Pagigation.Limit)
			suite.Equal(uint32(0), response.Pagigation.Offset)
			suite.Equal(uint32(1), response.Data[0].ID)
			suite.Equal("'Abaco' VW Bay Window: Westfalia Pop-top", response.Data[0].Name)
		}
	}
}

func (suite *RentalControllerTestSuite) TestListRentalsSuccessWithOffset() {
	req, _ := http.NewRequest("GET", "/rentals/?limit=1&offset=1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	if suite.Equal(http.StatusOK, w.Code) {
		var response testListResponse
		if suite.Nil(json.Unmarshal(w.Body.Bytes(), &response), "Should be able to unmarshal response") {
			suite.Equal(1, len(response.Data), "Should return a single result")
			suite.Equal(uint32(30), response.Pagigation.Count)
			suite.Equal(uint8(1), response.Pagigation.Limit)
			suite.Equal(uint32(1), response.Pagigation.Offset)
			suite.Equal(uint32(2), response.Data[0].ID)
			suite.Equal("Maupin: Vanagon Camper", response.Data[0].Name)
		}
	}
}

func (suite *RentalControllerTestSuite) TestListRentalsSuccessAllFilters() {
	req, _ := http.NewRequest("GET", "/rentals/?near=33.68,-117.82&price_min=9000&price_max=16000&ids=7,15&sort=price&limit=1&offset=0", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	if suite.Equal(http.StatusOK, w.Code) {
		var response testListResponse
		if suite.Nil(json.Unmarshal(w.Body.Bytes(), &response), "Should be able to unmarshal response") {
			suite.Equal(1, len(response.Data), "Should return a single result")
			suite.Equal(uint32(2), response.Pagigation.Count)
			suite.Equal(uint8(1), response.Pagigation.Limit)
			suite.Equal(uint32(0), response.Pagigation.Offset)
			suite.Equal(uint32(15), response.Data[0].ID)
			suite.Equal("AWESOME 1977 Volkswagen Westfalia camper", response.Data[0].Name)
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
		var response models.RentalResponse
		if suite.Nil(json.Unmarshal(w.Body.Bytes(), &response), "Should be able to unmarshal response") {
			suite.Equal(uint32(1), response.ID)
			suite.Equal(int64(16900), response.Price.Day)
			suite.Equal("Costa Mesa", response.Location.City)
		}
	}
}

func (suite *RentalControllerTestSuite) TestGetRentalIdNotFound() {
	req, _ := http.NewRequest("GET", "/rentals/100", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusNotFound, w.Code)
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
