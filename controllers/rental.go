package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/samuelg/rentals/logging"
	"github.com/samuelg/rentals/models"
)

type RentalController struct{}

// GET /rentals
func (u RentalController) List(c *gin.Context) {
	_, err := models.ParseQuery(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid filter", "error": err.Error()})
		c.Abort()
		return
	}

	// TODO: replace placeholder
	c.JSON(http.StatusOK, gin.H{"data": []string{}})
}

// GET /rentals/:rental_id
func (u RentalController) Get(c *gin.Context) {
	// id is an integer in the database, only needs int32
	rentalId, err := strconv.ParseInt(c.Param("rental_id"), 10, 32)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Invalid rental id: %s", c.Param("rental_id")))
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid rental id"})
		c.Abort()
		return
	}

	// TODO: replace placeholder
	c.JSON(http.StatusOK, gin.H{"id": rentalId})
}
