package rest

import (
	"github.com/gin-gonic/gin"
)

func ApplyRoutes(r *gin.Engine, h PlanetHandler) {
	planets := r.Group("/planets")
	{
		planets.GET("", h.GetPlanets)
		planets.POST("", h.CreatePlanet)
		planets.GET(":id", h.GetPlanetById)
		planets.GET("by-name", h.GetPlanetByName)
		planets.DELETE(":id", h.DeletePlanetById)

	}
}
