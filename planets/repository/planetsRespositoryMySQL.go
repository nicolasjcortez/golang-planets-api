package repository

import (
	"database/sql"
	"starwars/errs"
	"starwars/logger"
	"starwars/planets/domain"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type PlanetsRepositoryMySQL struct {
	Client *sqlx.DB
}

func (r PlanetsRepositoryMySQL) GetAllPlanets() ([]domain.Planet, *errs.AppError) {

	var planets []domain.Planet
	planetsSql := "select * from planets"
	err := r.Client.Select(&planets, planetsSql)

	if err != nil {
		logger.Error("Error while quering planets " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}
	return planets, nil
}

func (r PlanetsRepositoryMySQL) CreatePlanet(planet domain.PlanetCreationObj) (*string, *errs.AppError) {

	sqlInsert := "INSERT INTO planets (name, terrain, climate, qtd_films) VALUES (?,?,?,?)"
	result, err := r.Client.Exec(sqlInsert, planet.Name, planet.Terrain, planet.Climate, planet.QtdFilms)
	if err != nil {
		logger.Error("Error while inseting planets into db: " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	planetId, err := result.LastInsertId()
	if err != nil {
		logger.Error("Error while getting last inseterd id for new planet: " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	idStr := strconv.FormatInt(planetId, 10)

	return &idStr, nil

}

func (r PlanetsRepositoryMySQL) GetPlanetByName(name string) (*domain.Planet, *errs.AppError) {

	var planet domain.Planet
	planetsSql := "select * from planets where name = ?"
	err := r.Client.Get(&planet, planetsSql, name)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Planet not found")
		} else {
			logger.Error("Error while scaning planet " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}
	}
	return &planet, nil
}

func (r PlanetsRepositoryMySQL) GetPlanetById(id string) (*domain.Planet, *errs.AppError) {

	var planet domain.Planet
	planetsSql := "select * from planets where planet_id = ?"
	err := r.Client.Get(&planet, planetsSql, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Planet not found")
		} else {
			logger.Error("Error while scaning planet " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}
	}

	return &planet, nil
}

func (r PlanetsRepositoryMySQL) DeletePlanetById(id string) *errs.AppError {
	sqlDelete := "DELETE FROM planets WHERE planet_id = ?"
	_, err := r.Client.Exec(sqlDelete, id)
	if err != nil {
		logger.Error("Error while inseting planets into db: " + err.Error())
		return errs.NewUnexpectedError("unexpected database error")
	}
	return nil
}
