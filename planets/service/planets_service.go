package service

import (
	"errors"
	"starwars/domain"
	planets_repo "starwars/planets/repository/mongo"
)

func GetAllPlanets(host, database string) ([]domain.Planet, error) {
	return planets_repo.GetAllPlanets(host, database)
}

func getPlanetQtdFilms(planetName string) (int, error) {
	var err error
	var apiResult domain.ExternalPlanetsAPIResponse
	err = getPlanetExternalAPI(planetName, &apiResult)
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

func CreatePlanet(host, database string, planet domain.PlanetCreationRequest) (*domain.Planet, error) {
	_, err := planets_repo.GetPlanetByName(host, database, planet.Name)
	if err == nil {
		err := errors.New("Planet with this name already exists")
		return nil, err
	} else if err.Error() != "Planet with this name not found" {
		return nil, err
	}

	qtdFilms, err := getPlanetQtdFilms(planet.Name)
	if err != nil {
		err := errors.New("External Stawwars API not available or not working as expected")
		return nil, err
	}
	return planets_repo.CreatePlanet(host, database, planet, qtdFilms)
}

func GetPlanetByName(host, database, name string) (*domain.Planet, error) {
	return planets_repo.GetPlanetByName(host, database, name)
}

func GetPlanetById(host, database, id string) (*domain.Planet, error) {
	return planets_repo.GetPlanetById(host, database, id)
}

func DeletePlanetById(host, database, id string) error {
	_, err := planets_repo.GetPlanetById(host, database, id)
	if err != nil {
		return err
	}
	return planets_repo.DeletePlanetById(host, database, id)
}
