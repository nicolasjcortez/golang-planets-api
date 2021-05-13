package rest

import (
	"github.com/gin-gonic/gin"
)

func ApplyRoutes(r *gin.Engine) {
	planets := r.Group("/planets")
	{
		planets.GET("", GetPlanets)

	}
	planet := r.Group("/planet")
	{
		planet.POST("", CreatePlanet)
		planet.GET("by-id", GetPlanetById)
		planet.GET("by-name", GetPlanetByName)
		planet.DELETE("by-id", DeletePlanetById)

	}
}
