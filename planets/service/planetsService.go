package service

import (
	"errors"
	"starwars/planets/domain"
	"starwars/planets/repository"
)

type PlanetsService struct {
	Repo  repository.PlanetsRepository
	Swapi repository.PlanetsRepositorySwapi
}

func (s PlanetsService) GetAllPlanets() ([]domain.Planet, error) {
	return s.Repo.GetAllPlanets()
}

func (s PlanetsService) getPlanetQtdFilms(planetName string) (int, error) {
	var err error
	var apiResult domain.ExternalPlanetsAPIResponse
	err = s.Swapi.GetPlanetExternalAPI(planetName, &apiResult)
	if err != nil {
		return 0, err
	}
	if apiResult.Count == 0 || len(apiResult.Results) == 0 {
		return 0, nil
	}

	if apiResult.Results[0].Name != planetName {
		err := errors.New("External Stawwars API not working as expected")
		return 0, err
	}
	qtdFilms := len(apiResult.Results[0].Films)
	return qtdFilms, nil
}

func (s PlanetsService) CreatePlanet(planet domain.PlanetCreationRequest) (*domain.Planet, error) {
	_, err := s.Repo.GetPlanetByName(planet.Name)
	if err == nil {
		err := errors.New("Planet with this name already exists")
		return nil, err
	} else if err.Error() != "Planet with this name not found" {
		return nil, err
	}

	qtdFilms, err := s.getPlanetQtdFilms(planet.Name)
	if err != nil {
		err := errors.New("External Stawwars API not available or not working as expected")
		return nil, err
	}
	return s.Repo.CreatePlanet(planet, qtdFilms)
}

func (s PlanetsService) GetPlanetByName(name string) (*domain.Planet, error) {
	return s.Repo.GetPlanetByName(name)
}

func (s PlanetsService) GetPlanetById(id string) (*domain.Planet, error) {
	return s.Repo.GetPlanetById(id)
}

func (s PlanetsService) DeletePlanetById(id string) error {
	_, err := s.Repo.GetPlanetById(id)
	if err != nil {
		return err
	}
	return s.Repo.DeletePlanetById(id)
}
