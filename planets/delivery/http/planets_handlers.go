package rest

import (
	"errors"
	"net/http"

	general_domain "starwars/domain"
	"starwars/planets/domain"
	planets_service "starwars/planets/service"

	"github.com/gin-gonic/gin"
)

type PlanetHandler struct {
	PlanetsService planets_service.PlanetsService
}

// @Accept  json
// @Produce  json
// @Success 200 {object} []domain.Planet
// @Failure 500 {object} domain.GinError
// @Router /planets [get]
func (h PlanetHandler) GetPlanets(c *gin.Context) {

	result, err := h.PlanetsService.GetAllPlanets()
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
// @Param _id query string true "Planet Database id"
// @Failure 400 {object} domain.GinError
// @Failure 404 {object} domain.GinError
// @Failure 500 {object} domain.GinError
// @Router /planet/by-id [get]
func (h PlanetHandler) GetPlanetById(c *gin.Context) {

	id := c.Query("_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, c.Error(errors.New("Missing query parameter: _id")))
		return
	}

	result, err := h.PlanetsService.GetPlanetById(id)

	if err != nil {
		switch err.Error() {
		case "Planet with this id not found":
			c.JSON(http.StatusNotFound, c.Error(err))
			return
		case "Error converting query parameter id to Database id format":
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
func (h PlanetHandler) GetPlanetByName(c *gin.Context) {

	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, c.Error(errors.New("Missing query parameter: name")))
		return
	}

	result, err := h.PlanetsService.GetPlanetByName(name)
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
func (h PlanetHandler) CreatePlanet(c *gin.Context) {

	var planet domain.PlanetCreationRequest
	err := c.BindJSON(&planet)
	if err != nil {
		c.JSON(http.StatusBadRequest, c.Error(err))
		return
	}

	result, err := h.PlanetsService.CreatePlanet(planet)
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
// @Param _id query string true "Planet Database id"
// @Failure 400 {object} domain.GinError
// @Failure 404 {object} domain.GinError
// @Failure 500 {object} domain.GinError
// @Router /planet/by-id [delete]
func (h PlanetHandler) DeletePlanetById(c *gin.Context) {

	id := c.Query("_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, c.Error(errors.New("Missing query parameter: _id")))
		return
	}

	err := h.PlanetsService.DeletePlanetById(id)

	if err != nil {
		switch err.Error() {
		case "Planet with this id not found":
			c.JSON(http.StatusNotFound, c.Error(err))
			return
		case "Error converting query parameter id to Database id format":
			c.JSON(http.StatusBadRequest, c.Error(err))
			return
		default:
			c.JSON(500, c.Error(err))
			return
		}
	} else {
		success := general_domain.SuccessResponse{
			Result: "Deleted Successfully",
		}
		c.JSON(http.StatusOK, success)
		return
	}
}
