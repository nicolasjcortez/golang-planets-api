package repository

import "starwars/planets/domain"

type PlanetsRepository interface {
	GetAllPlanets() ([]domain.Planet, error)
	CreatePlanet(planetRequest domain.PlanetCreationRequest, qtdFilms int) (*domain.Planet, error)
	GetPlanetByName(name string) (*domain.Planet, error)
	GetPlanetById(id string) (*domain.Planet, error)
	DeletePlanetById(id string) error
}
