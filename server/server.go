package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/samuelg/rentals/config"
	"github.com/samuelg/rentals/controllers"
	log "github.com/samuelg/rentals/logging"
)

// Create router with routes for rentals
func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// TODO: in production trust load balancer in front of API
	router.SetTrustedProxies(nil)

	rentalGroup := router.Group("rentals")
	{
		rentals := new(controllers.RentalController)
		rentalGroup.GET("/", rentals.List)
		rentalGroup.GET("/:rental_id", rentals.Get)
	}

	log.Log.Info("Router created")

	return router
}

func Init() {
	r := NewRouter()

	listenAddress := fmt.Sprintf("%s:%d", config.GetConfig().Host, config.GetConfig().Port)
	log.Log.Info(fmt.Sprintf("Listening on %s", listenAddress))
	log.Log.Fatal(r.Run(listenAddress))
}
