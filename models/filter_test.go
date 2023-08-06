package models

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/samuelg/rentals/config"
	"github.com/samuelg/rentals/db"
	log "github.com/samuelg/rentals/logging"
	"github.com/stretchr/testify/suite"
)

// Test suite for the Rental controller
type FilterModelTestSuite struct {
	suite.Suite
	config *config.Config
}

func (suite *FilterModelTestSuite) SetupSuite() {
	config.Init("test")
	log.Init("FATAL", config.GetConfig().AppVersion)
	db.Init()
	suite.config = config.GetConfig()
}

// Mock a gin request with a query to test our filter query parser
func mockQuery(q url.Values) *gin.Context {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := &http.Request{
		URL: &url.URL{},
	}

	req.URL.RawQuery = q.Encode()
	c.Request = req

	return c
}

// filter.ParseQuery tests
func (suite *FilterModelTestSuite) TestParseQuerySuccessAllFilters() {
	q := url.Values{}
	q.Set("price_min", "1000")
	q.Set("price_max", "2000")
	q.Set("limit", "5")
	q.Set("offset", "0")
	q.Set("near", "33.68,-117.82")
	q.Set("ids", "1,2,3")
	q.Set("sort", "price")
	c := mockQuery(q)

	filter, err := ParseQuery(c)

	if suite.Nil(err, "Should not result in an error") {
		suite.Equal(int64(1000), *filter.PriceMin)
		suite.Equal(int64(2000), *filter.PriceMax)
		suite.Equal(uint8(5), filter.Limit)
		suite.Equal(uint32(0), filter.Offset)
		suite.Equal(float32(33.68), filter.Near[0])
		suite.Equal(float32(-117.82), filter.Near[1])
		suite.Equal(uint32(1), filter.Ids[0])
		suite.Equal(uint32(2), filter.Ids[1])
		suite.Equal(uint32(3), filter.Ids[2])
		suite.Equal("price", filter.Sort)
	}
}

func (suite *FilterModelTestSuite) TestParseQuerySuccessDefaults() {
	q := url.Values{}
	c := mockQuery(q)

	filter, err := ParseQuery(c)

	if suite.Nil(err, "Should not result in an error") {
		suite.Nil(filter.PriceMin, "Should not be assigned")
		suite.Nil(filter.PriceMax, "Should not be assigned")
		// default limit in tests config is set to 1
		suite.Equal(uint8(1), filter.Limit)
		suite.Equal(uint32(0), filter.Offset)
		suite.Equal(len(filter.Near), 0)
		suite.Equal(len(filter.Ids), 0)
		suite.Equal("", filter.Sort)
	}
}

func (suite *FilterModelTestSuite) TestParseQueryInvalidPriceMin() {
	q := url.Values{}
	q.Set("price_min", "abc")
	c := mockQuery(q)

	_, err := ParseQuery(c)

	if suite.NotNil(err, "Should result in an error") {
		suite.Equal("Invalid price_min", err.Error())
	}
}

func (suite *FilterModelTestSuite) TestParseQueryInvalidPriceMax() {
	q := url.Values{}
	q.Set("price_max", "abc")
	c := mockQuery(q)

	_, err := ParseQuery(c)

	if suite.NotNil(err, "Should result in an error") {
		suite.Equal("Invalid price_max", err.Error())
	}
}

func (suite *FilterModelTestSuite) TestParseQueryInvalidLimit() {
	q := url.Values{}
	q.Set("limit", "abc")
	c := mockQuery(q)

	_, err := ParseQuery(c)

	if suite.NotNil(err, "Should result in an error") {
		suite.Equal("Invalid limit", err.Error())
	}
}

func (suite *FilterModelTestSuite) TestParseQueryLimitTooLarge() {
	q := url.Values{}
	q.Set("limit", "101")
	c := mockQuery(q)

	_, err := ParseQuery(c)

	if suite.NotNil(err, "Should result in an error") {
		suite.Equal("Limit is too large", err.Error())
	}
}

func (suite *FilterModelTestSuite) TestParseQueryInvalidOffset() {
	q := url.Values{}
	q.Set("offset", "abc")
	c := mockQuery(q)

	_, err := ParseQuery(c)

	if suite.NotNil(err, "Should result in an error") {
		suite.Equal("Invalid offset", err.Error())
	}
}

func (suite *FilterModelTestSuite) TestParseQueryInvalidSort() {
	q := url.Values{}
	q.Set("sort", "abc")
	c := mockQuery(q)

	_, err := ParseQuery(c)

	if suite.NotNil(err, "Should result in an error") {
		suite.Equal("Invalid sort", err.Error())
	}
}

func (suite *FilterModelTestSuite) TestParseQueryInvalidIds() {
	q := url.Values{}
	q.Set("ids", "abc")
	c := mockQuery(q)

	_, err := ParseQuery(c)

	if suite.NotNil(err, "Should result in an error") {
		suite.Equal("Invalid id in ids", err.Error())
	}
}

func (suite *FilterModelTestSuite) TestParseQueryTooFewNear() {
	q := url.Values{}
	q.Set("near", "36.23")
	c := mockQuery(q)

	_, err := ParseQuery(c)

	if suite.NotNil(err, "Should result in an error") {
		suite.Equal("Invalid value for near", err.Error())
	}
}

func (suite *FilterModelTestSuite) TestParseQueryInvalidLatitude() {
	q := url.Values{}
	q.Set("near", "abc,-117.12")
	c := mockQuery(q)

	_, err := ParseQuery(c)

	if suite.NotNil(err, "Should result in an error") {
		suite.Equal("Invalid latitude", err.Error())
	}
}

func (suite *FilterModelTestSuite) TestParseQueryInvalidLongitude() {
	q := url.Values{}
	q.Set("near", "36.23,abc")
	c := mockQuery(q)

	_, err := ParseQuery(c)

	if suite.NotNil(err, "Should result in an error") {
		suite.Equal("Invalid longitude", err.Error())
	}
}

func (suite *FilterModelTestSuite) TestParseQueryLatitudeTooLarge() {
	q := url.Values{}
	q.Set("near", "190.0,-117.12")
	c := mockQuery(q)

	_, err := ParseQuery(c)

	if suite.NotNil(err, "Should result in an error") {
		suite.Equal("Invalid latitude", err.Error())
	}
}

func (suite *FilterModelTestSuite) TestParseQueryLatitudeTooSmall() {
	q := url.Values{}
	q.Set("near", "-190.0,-117.12")
	c := mockQuery(q)

	_, err := ParseQuery(c)

	if suite.NotNil(err, "Should result in an error") {
		suite.Equal("Invalid latitude", err.Error())
	}
}

func (suite *FilterModelTestSuite) TestParseQueryLongitudeTooLarge() {
	q := url.Values{}
	q.Set("near", "36.23,190.0")
	c := mockQuery(q)

	_, err := ParseQuery(c)

	if suite.NotNil(err, "Should result in an error") {
		suite.Equal("Invalid longitude", err.Error())
	}
}

func (suite *FilterModelTestSuite) TestParseQueryLongitudeTooSmall() {
	q := url.Values{}
	q.Set("near", "36.23,-190.0")
	c := mockQuery(q)

	_, err := ParseQuery(c)

	if suite.NotNil(err, "Should result in an error") {
		suite.Equal("Invalid longitude", err.Error())
	}
}

// filter.Find tests
func (suite *FilterModelTestSuite) TestFindSuccessAllFilters() {
	priceMin := int64(9000)
	priceMax := int64(16000)

	filter := &Filter{
		PriceMin: &priceMin,
		PriceMax: &priceMax,
		Limit:    uint8(1),
		Offset:   uint32(0),
		// lat / lng near Costa Mesa
		Near: []float32{float32(33.68), float32(-117.82)},
		// we would expect to see 7 first but the sort on price will flip the ids
		Ids:  []uint32{7, 15},
		Sort: "price",
	}

	rentals, count, err := filter.Find()

	if suite.Nil(err, "Should not lead to an error") {
		suite.Equal(uint32(2), count)
		suite.Equal(len(rentals), 1)
		// will get 15 first as it sorts on price
		suite.Equal(uint32(15), rentals[0].ID)
		suite.Equal("AWESOME 1977 Volkswagen Westfalia camper", rentals[0].Name)
	}
}

func (suite *FilterModelTestSuite) TestFindSuccessNoFilters() {
	// ParseQuery will set these by default for tests due to default test limit
	filter := &Filter{Limit: 1, Offset: 0, Sort: "id"}

	rentals, count, err := filter.Find()

	if suite.Nil(err, "Should not lead to an error") {
		suite.Equal(uint32(30), count)
		suite.Equal(len(rentals), 1)
		// ParseQuery sorts on id by default
		suite.Equal(uint32(1), rentals[0].ID)
		suite.Equal("'Abaco' VW Bay Window: Westfalia Pop-top", rentals[0].Name)
	}
}

func TestFilterModelTestSuite(t *testing.T) {
	suite.Run(t, new(FilterModelTestSuite))
}
