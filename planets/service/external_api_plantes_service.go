package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"starwars/domain"
)

func getParsedJsonExternalPlanetAPI(url string, target *domain.ExternalPlanetsAPIResponse) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(&target)
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

func getPlanetExternalAPI(planetName string, apiResult *domain.ExternalPlanetsAPIResponse) error {
	externalPlanetsURL, err := buildExternalPlanetURL(planetName)
	if err != nil {
		return err
	}
	fmt.Println(externalPlanetsURL)

	return getParsedJsonExternalPlanetAPI(externalPlanetsURL, apiResult)
}
