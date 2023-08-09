package metrics

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/samuelg/rentals/config"
	log "github.com/samuelg/rentals/logging"
	"github.com/stretchr/testify/suite"
)

// Test suite for the Receipt model
type MetricsTestSuite struct {
	suite.Suite
	config *config.Config
	router *gin.Engine
}

func setupRouter() *gin.Engine {
	// don't log API calls
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// init metrics
	Init(router)

	return router
}

func (suite *MetricsTestSuite) SetupSuite() {
	config.Init("test")
	log.Init("FATAL", config.GetConfig().AppVersion)
	suite.config = config.GetConfig()
	suite.router = setupRouter()
}

func (suite *MetricsTestSuite) TestDefaultMetrics() {
	req, _ := http.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	if suite.Equal(http.StatusOK, w.Code) {
		res, err := ioutil.ReadAll(w.Result().Body)
		if suite.Nil(err, "Error should be nil") {
			body := string(res)
			suite.True(
				strings.Contains(body, "promhttp_metric_handler_requests_total"),
				"Response should contain default metrics",
			)
		}
	}
}

func TestMetricsTestSuite(t *testing.T) {
	suite.Run(t, new(MetricsTestSuite))
}
