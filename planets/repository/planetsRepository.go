package repository

import (
	"starwars/errs"
	"starwars/planets/domain"
)

type PlanetsRepository interface {
	GetAllPlanets() ([]domain.Planet, *errs.AppError)
	CreatePlanet(planetRequest domain.PlanetCreationRequest, qtdFilms int) (*domain.Planet, *errs.AppError)
	GetPlanetByName(name string) (*domain.Planet, *errs.AppError)
	GetPlanetById(id string) (*domain.Planet, *errs.AppError)
	DeletePlanetById(id string) *errs.AppError
}
