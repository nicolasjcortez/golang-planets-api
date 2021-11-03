package repository

import (
	"starwars/errs"
	"starwars/planets/domain"
)

type PlanetsRepository interface {
	GetAllPlanets() ([]domain.Planet, *errs.AppError)
	CreatePlanet(planet domain.PlanetCreationObj, qtdFilms int) (*string, *errs.AppError)
	GetPlanetByName(name string) (*domain.Planet, *errs.AppError)
	GetPlanetById(id string) (*domain.Planet, *errs.AppError)
	DeletePlanetById(id string) *errs.AppError
}
