package domain

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Planet struct {
	ID       IDType `json:"_id" bson:"_id,omitempty" binding:"required"`
	Name     string `json:"name" bson:"name" binding:"required"`
	Climate  string `json:"climate" bson:"climate" binding:"required"`
	Terrain  string `json:"terrain" bson:"terrain" binding:"required"`
	QtdFilms int    `json:"qtd_films" bson:"qtd_films" binding:"required"`
}

type IDType string

func (id IDType) MarshalBSONValue() (bsontype.Type, []byte, error) {
	p, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return bsontype.Null, nil, err
	}

	return bson.MarshalValue(p)
}

type PlanetCreationRequest struct {
	Name    string `json:"name" bson:"name" binding:"required"`
	Climate string `json:"climate" bson:"climate" binding:"required"`
	Terrain string `json:"terrain" bson:"terrain" binding:"required"`
}

type PlanetCreationObj struct {
	Name     string `json:"name" bson:"name"`
	Climate  string `json:"climate" bson:"climate"`
	Terrain  string `json:"terrain" bson:"terrain"`
	QtdFilms int    `json:"qtd_films" bson:"qtd_films"`
}

type ExternalPlanetsAPIResponse struct {
	Count   int                 `json:"count" bson:"count"`
	Results []ExternalAPIPlanet `json:"results" bson:"results"`
}

type ExternalAPIPlanet struct {
	Name  string   `json:"name" bson:"name"`
	Films []string `json:"films" bson:"films"`
}
