package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/samuelg/rentals/db"
	log "github.com/samuelg/rentals/logging"
	"github.com/samuelg/rentals/models"
	"gorm.io/gorm"
)

type RentalController struct{}

// GET /rentals
func (u RentalController) List(c *gin.Context) {
	filter, err := models.ParseQuery(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid filter", "error": err.Error()})
		c.Abort()
		return
	}

	rentals, err := filter.Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong", "error": err.Error()})
		c.Abort()
		return
	}

	// TODO: return {pagination, data}
	c.JSON(http.StatusOK, rentals)
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

	var rental models.Rental
	if result := db.DB.Joins("User").First(&rental, rentalId); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Rental not found"})
			c.Abort()
			return
		}
	}

	c.JSON(http.StatusOK, rental)
}
