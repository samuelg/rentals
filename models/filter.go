package models

import (
	"errors"
	"fmt"
	"golang.org/x/exp/slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/samuelg/rentals/config"
	"github.com/samuelg/rentals/db"
	log "github.com/samuelg/rentals/logging"
)

// Represents a filter on a list of rentals
type Filter struct {
	// All query params are optional
	PriceMin *int64
	PriceMax *int64
	Limit    uint8 // don't allow a large limit value
	Offset   int32
	Ids      []int32
	Near     []float32
	Sort     string
}

func ParseQuery(c *gin.Context) (*Filter, error) {
	filter := new(Filter)
	// store error messages as we discover them
	validationErrors := make([]string, 0)

	priceMinRaw := c.Query("price_min")
	if priceMinRaw != "" {
		price, err := strconv.ParseInt(priceMinRaw, 10, 64)
		if err != nil {
			log.Log.Trace(fmt.Sprintf("Invalid price_min: %s", priceMinRaw))
			validationErrors = append(validationErrors, "Invalid price_min")
		} else {
			filter.PriceMin = &price
		}
	}

	priceMaxRaw := c.Query("price_max")
	if priceMaxRaw != "" {
		price, err := strconv.ParseInt(priceMaxRaw, 10, 64)
		if err != nil {
			log.Log.Trace(fmt.Sprintf("Invalid price_max: %s", priceMaxRaw))
			validationErrors = append(validationErrors, "Invalid price_max")
		} else {
			filter.PriceMax = &price
		}
	}

	limitRaw := c.Query("limit")
	if limitRaw != "" {
		limit, err := strconv.ParseInt(limitRaw, 10, 8)
		if err != nil {
			log.Log.Trace(fmt.Sprintf("Invalid limit: %s", limitRaw))
			validationErrors = append(validationErrors, "Invalid limit")
		} else {
			// don't allow a large value for limit
			if limit > 100 {
				log.Log.Trace(fmt.Sprintf("Limit is too large: %d", limit))
				validationErrors = append(validationErrors, "Limit is too large")
			} else {
				filter.Limit = uint8(limit)
			}
		}
	}
	// use default limit if not provided (offset will default to 0)
	if filter.Limit == 0 {
		filter.Limit = config.GetConfig().DefaultApiLimit
	}

	offsetRaw := c.Query("offset")
	if offsetRaw != "" {
		offset, err := strconv.ParseInt(offsetRaw, 10, 32)
		if err != nil {
			log.Log.Trace(fmt.Sprintf("Invalid offset: %s", offsetRaw))
			validationErrors = append(validationErrors, "Invalid offset")
		} else {
			filter.Offset = int32(offset)
		}
	}

	validSorts := []string{"", "name", "type", "sleeps", "price", "city", "state", "country", "make", "model", "year", "length", "created", "updated"}
	sort := c.Query("sort")
	if !slices.Contains(validSorts, sort) {
		validationErrors = append(validationErrors, "Invalid sort")
	} else {
		filter.Sort = c.Query("sort")
	}

	idsRaw := c.Query("ids")
	// parse csv ids
	if idsRaw != "" {
		ids := strings.Split(idsRaw, ",")
		for _, idRaw := range ids {
			id, err := strconv.ParseInt(idRaw, 10, 32)
			if err != nil {
				log.Log.Trace(fmt.Sprintf("Invalid id: %s", idRaw))
				validationErrors = append(validationErrors, "Invalid id in ids")
				break
			}
			filter.Ids = append(filter.Ids, int32(id))
		}
	}

	nearRaw := c.Query("near")
	// parse csv near lat/long
	if nearRaw != "" {
		latLong := strings.Split(nearRaw, ",")
		if len(latLong) != 2 {
			log.Log.Trace(fmt.Sprintf("Too many values for lat/long: %d", len(latLong)))
			validationErrors = append(validationErrors, "Invalid value for near")
		} else {
			lat, latErr := strconv.ParseFloat(latLong[0], 32)
			long, longErr := strconv.ParseFloat(latLong[1], 32)

			if latErr != nil {
				log.Log.Trace(fmt.Sprintf("Invalid lat: %s", latLong[0]))
				validationErrors = append(validationErrors, "Invalid latitude")
			}
			if longErr != nil {
				log.Log.Trace(fmt.Sprintf("Invalid long: %s", latLong[1]))
				validationErrors = append(validationErrors, "Invalid longitude")
			}

			if latErr == nil && longErr == nil {
				// validate lat / long values
				validLat := lat <= 90 && lat >= -90
				validLong := long <= 180 && long >= -180
				if !validLat {
					log.Log.Trace(fmt.Sprintf("Invalid lat: %.2f", lat))
					validationErrors = append(validationErrors, "Invalid latitude")
				}
				if !validLong {
					log.Log.Trace(fmt.Sprintf("Invalid long: %.2f", long))
					validationErrors = append(validationErrors, "Invalid longitude")
				}

				if validLat && validLong {
					filter.Near = []float32{float32(lat), float32(long)}
				}
			}
		}
	}

	if len(validationErrors) > 0 {
		return nil, errors.New(strings.Join(validationErrors, "\n"))
	}

	return filter, nil
}

// Find rentals using the provided filter
// TODO: return count as well within transaction
func (filter *Filter) Find() ([]Rental, error) {
	var rentals []Rental
	// Default query
	query := db.DB.Joins("User")

	// Minimum price
	if filter.PriceMin != nil {
		query = query.Where("price_per_day >= ?", *filter.PriceMin)
	}

	// Maximum price
	if filter.PriceMax != nil {
		query = query.Where("price_per_day <= ?", *filter.PriceMax)
	}

	// IDs
	if len(filter.Ids) != 0 {
		// IN clause
		query = query.Where(filter.Ids)
	}

	// Near
	if len(filter.Near) == 2 {
		lat := filter.Near[0]
		lng := filter.Near[1]
		// Use geography to calculate in meters, then convert from miles
		query = query.Where(
			"ST_DWITHIN(ST_SETSRID(st_makepoint(lng, lat), 4326)::geography, st_setsrid(st_makepoint(?, ?), 4326)::geography, 100 * 1609.34)",
			lng,
			lat,
		)
	}

	// Limit sort to known values
	var sort string
	switch filter.Sort {
	case "name":
		sort = "name"
	case "type":
		sort = "type"
	case "sleeps":
		sort = "sleeps"
	case "price":
		sort = "price_per_day"
	case "city":
		sort = "home_city"
	case "state":
		sort = "home_state"
	case "country":
		sort = "home_country"
	case "make":
		sort = "vehicle_make"
	case "model":
		sort = "vehicle_model"
	case "year":
		sort = "vehicle_year"
	case "length":
		sort = "vehicle_length"
	case "created":
		sort = "created"
	case "updated":
		sort = "updated"
	default:
		// Default sort
		sort = "id"
	}

	// Apply limit and offset
	query = query.Limit(int(filter.Limit)).Offset(int(filter.Offset))
	// Apply sort
	query = query.Order(sort)

	if result := query.Find(&rentals); result.Error != nil {
		return nil, result.Error
	}

	return rentals, nil
}
