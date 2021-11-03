package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"starwars/errs"
	"starwars/planets/domain"
)

type PlanetsRepositorySwapi interface {
	GetPlanetExternalAPI(string) (*domain.ExternalPlanetsAPIResponse, *errs.AppError)
}

type RemotePlanetsRespositorySwapi struct {
}

func getParsedJsonExternalPlanetAPI(url string, target *domain.ExternalPlanetsAPIResponse) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&target)
	return err
}

func buildExternalPlanetURL(planetName string) (string, error) {
	base, err := url.Parse("https://swapi.dev/api/planets")
	if err != nil {
		return "", err
	}

	// Query params
	params := url.Values{}
	params.Add("search", planetName)
	base.RawQuery = params.Encode()

	// fmt.Printf("Encoded URL is %q\n", base.String())
	return base.String(), nil
}

func (repo RemotePlanetsRespositorySwapi) GetPlanetExternalAPI(planetName string) (*domain.ExternalPlanetsAPIResponse, *errs.AppError) {
	externalPlanetsURL, err := buildExternalPlanetURL(planetName)
	if err != nil {
		return nil, errs.NewUnexpectedError("External Stawwars API not working as expected")
	}
	fmt.Println(externalPlanetsURL)
	var apiResult domain.ExternalPlanetsAPIResponse
	err = getParsedJsonExternalPlanetAPI(externalPlanetsURL, &apiResult)
	if err != nil {
		return nil, errs.NewUnexpectedError("External Stawwars API not working as expected")
	}
	return &apiResult, nil
}
