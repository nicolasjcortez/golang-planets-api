package service

import (
	"net/http"
	"starwars/errs"
	"starwars/planets/domain"
	"starwars/planets/dto"
	"starwars/planets/repository"
	"strconv"
)

type PlanetsService struct {
	Repo  repository.PlanetsRepository
	Swapi repository.PlanetsRepositorySwapi
}

func (s PlanetsService) GetAllPlanets() ([]domain.Planet, *errs.AppError) {
	return s.Repo.GetAllPlanets()
}

func (s PlanetsService) getPlanetQtdFilms(planetName string) (int, *errs.AppError) {
	apiResult, err := s.Swapi.GetPlanetExternalAPI(planetName)
	if err != nil {
		return 0, err
	}
	if apiResult.Count == 0 || len(apiResult.Results) == 0 {
		return 0, nil
	}

	if apiResult.Results[0].Name != planetName {
		return 0, err
	}
	qtdFilms := len(apiResult.Results[0].Films)
	return qtdFilms, nil
}

func (s PlanetsService) CreatePlanet(planetRequest dto.PlanetCreationRequest) (*domain.Planet, *errs.AppError) {
	_, err := s.Repo.GetPlanetByName(planetRequest.Name)
	if err == nil {
		return nil, errs.NewConflictError("Planet with this name already exists")
	} else if err.Code != http.StatusNotFound {
		return nil, err
	}

	qtdFilms, err := s.getPlanetQtdFilms(planetRequest.Name)
	if err != nil {
		return nil, err
	}

	planetDB := domain.PlanetCreationObj{
		Name:     planetRequest.Name,
		Climate:  planetRequest.Climate,
		Terrain:  planetRequest.Terrain,
		QtdFilms: qtdFilms,
	}

	id, err := s.Repo.CreatePlanet(planetDB)
	if err != nil {
		return nil, err
	}

	// sql version
	idInt, _ := strconv.Atoi(*id)

	responsePlanet := domain.Planet{
		// ID:       domain.IDType(*id), //mongo version
		ID:       idInt, //sql version
		Name:     planetDB.Name,
		Climate:  planetDB.Climate,
		Terrain:  planetDB.Terrain,
		QtdFilms: qtdFilms,
	}

	return &responsePlanet, nil
}

func (s PlanetsService) GetPlanetByName(name string) (*domain.Planet, *errs.AppError) {
	if name == "" {
		return nil, errs.NewBadRequestError("Missing query parameter: name")
	}
	return s.Repo.GetPlanetByName(name)
}

func (s PlanetsService) GetPlanetById(id string) (*domain.Planet, *errs.AppError) {
	if id == "" {
		return nil, errs.NewBadRequestError("Missing query parameter: _id")
	}
	return s.Repo.GetPlanetById(id)
}

func (s PlanetsService) DeletePlanetById(id string) *errs.AppError {
	_, err := s.Repo.GetPlanetById(id)
	if err != nil {
		return err
	}
	return s.Repo.DeletePlanetById(id)
}
