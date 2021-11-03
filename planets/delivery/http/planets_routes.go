package rest

import (
	"github.com/gin-gonic/gin"
)

func ApplyRoutes(r *gin.Engine, h PlanetHandler) {
	planets := r.Group("/planets")
	{
		planets.GET("", h.GetPlanets)

	}
	planet := r.Group("/planet")
	{
		planet.POST("", h.CreatePlanet)
		planet.GET("by-id", h.GetPlanetById)
		planet.GET("by-name", h.GetPlanetByName)
		planet.DELETE("by-id", h.DeletePlanetById)

	}
}
