package rest

import (
	"errors"
	"net/http"

	"starwars/domain"
	planets_service "starwars/planets/service"

	"github.com/gin-gonic/gin"
)

// @Accept  json
// @Produce  json
// @Success 200 {object} []domain.Planet
// @Failure 500 {object} domain.GinError
// @Router /planets [get]
func GetPlanets(c *gin.Context) {
	host := c.MustGet("host").(string)
	database := c.MustGet("db").(string)

	result, err := planets_service.GetAllPlanets(
		host,
		database)
	if err != nil {
		c.JSON(500, c.Error(err))
		return
	} else {
		c.JSON(http.StatusOK, result)
		return
	}
}

// @Accept  json
// @Produce  json
// @Success 200 {object} domain.Planet
// @Param _id query string true "Planet database id"
// @Failure 400 {object} domain.GinError
// @Failure 404 {object} domain.GinError
// @Failure 500 {object} domain.GinError
// @Router /planet/by-id [get]
func GetPlanetById(c *gin.Context) {
	host := c.MustGet("host").(string)
	database := c.MustGet("db").(string)
	id := c.Query("_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, c.Error(errors.New("Missing query parameter: _id")))
		return
	}

	result, err := planets_service.GetPlanetById(
		host,
		database,
		id)

	if err != nil {
		switch err.Error() {
		case "Planet with this id not found":
			c.JSON(http.StatusNotFound, c.Error(err))
			return
		case "Error converting query parameter id to database id format":
			c.JSON(http.StatusBadRequest, c.Error(err))
			return
		default:
			c.JSON(500, c.Error(err))
			return
		}
	} else {
		c.JSON(http.StatusOK, result)
		return
	}
}

// @Accept  json
// @Produce  json
// @Success 200 {object} domain.Planet
// @Param name query string true "Planet name"
// @Failure 400 {object} domain.GinError
// @Failure 404 {object} domain.GinError
// @Failure 500 {object} domain.GinError
// @Router /planet/by-name [get]
func GetPlanetByName(c *gin.Context) {
	host := c.MustGet("host").(string)
	database := c.MustGet("db").(string)
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, c.Error(errors.New("Missing query parameter: name")))
		return
	}

	result, err := planets_service.GetPlanetByName(
		host,
		database,
		name)
	if err != nil {
		switch err.Error() {
		case "Planet with this name not found":
			c.JSON(http.StatusNotFound, c.Error(err))
			return
		default:
			c.JSON(500, c.Error(err))
			return
		}
	} else {
		c.JSON(http.StatusOK, result)
		return
	}
}

// @Accept  json
// @Produce  json
// @Param planet body domain.PlanetCreationRequest true "Planet Data"
// @Success 201 {object} domain.Planet
// @Failure 400 {object} domain.GinError
// @Failure 500 {object} domain.GinError
// @Router /planet [post]
func CreatePlanet(c *gin.Context) {
	host := c.MustGet("host").(string)
	database := c.MustGet("db").(string)

	var planet domain.PlanetCreationRequest
	err := c.BindJSON(&planet)
	if err != nil {
		c.JSON(http.StatusBadRequest, c.Error(err))
		return
	}

	result, err := planets_service.CreatePlanet(
		host,
		database,
		planet)
	if err != nil {
		switch err.Error() {
		case "Planet with this name not found":
			c.JSON(http.StatusNotFound, c.Error(err))
			return
		case "Planet with this name already exists":
			c.JSON(422, c.Error(err))
			return
		default:
			c.JSON(500, c.Error(err))
			return
		}
	} else {
		c.JSON(http.StatusCreated, result)
		return
	}
}

// @Accept  json
// @Produce  json
// @Success 200 {object} domain.SuccessResponse
// @Param _id query string true "Planet database id"
// @Failure 400 {object} domain.GinError
// @Failure 404 {object} domain.GinError
// @Failure 500 {object} domain.GinError
// @Router /planet/by-id [delete]
func DeletePlanetById(c *gin.Context) {
	host := c.MustGet("host").(string)
	database := c.MustGet("db").(string)
	id := c.Query("_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, c.Error(errors.New("Missing query parameter: _id")))
		return
	}

	err := planets_service.DeletePlanetById(
		host,
		database,
		id)

	if err != nil {
		switch err.Error() {
		case "Planet with this id not found":
			c.JSON(http.StatusNotFound, c.Error(err))
			return
		case "Error converting query parameter id to database id format":
			c.JSON(http.StatusBadRequest, c.Error(err))
			return
		default:
			c.JSON(500, c.Error(err))
			return
		}
	} else {
		success := domain.SuccessResponse{
			Result: "Deleted Successfully",
		}
		c.JSON(http.StatusOK, success)
		return
	}
}
