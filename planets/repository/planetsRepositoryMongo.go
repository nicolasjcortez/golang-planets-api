package repository

import (
	"context"
	"fmt"
	"log"
	"starwars/errs"
	"starwars/logger"
	"starwars/planets/domain"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var once sync.Once

type MongoConnection struct {
	Conn *mongo.Client
}

var instance *MongoConnection

type PlanetsRepositoryMongo struct {
	Host     string
	Database string
}

func Connect(host string) *MongoConnection {
	once.Do(func() {
		fmt.Print("SINGLETON CONNECTION CREATED")
		ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(host).
			SetMaxPoolSize(250).
			SetConnectTimeout(160*time.Second).
			SetServerSelectionTimeout(160*time.Second).
			SetSocketTimeout(160*time.Second))
		if err != nil {
			// maybe refactor later, fatal shutdown server
			log.Fatal(err)
		}
		instance = &MongoConnection{Conn: client}
	})
	return instance
}

func (r PlanetsRepositoryMongo) GetAllPlanets() ([]domain.Planet, *errs.AppError) {
	var err error
	var ctx, _ = context.WithTimeout(context.Background(), 40*time.Second)
	var planets []domain.Planet

	// create the context and connection (once)
	client := Connect(r.Host)

	// select the db and collection
	collection := client.Conn.Database(r.Database).Collection("planets")

	// set options
	options := options.Find()

	cur, err := collection.Find(ctx, bson.M{}, options)
	defer cur.Close(ctx)
	if err != nil {
		log.Print(err)
		return nil, errs.NewUnexpectedError("Unexpected database error finding planets")
	}

	if err := cur.All(ctx, &planets); err != nil {
		return nil, errs.NewUnexpectedError("Unexpected database error finding planets")
	}

	if planets == nil {
		planets = make([]domain.Planet, 0)
		return planets, nil
	}

	return planets, nil
}

func (r PlanetsRepositoryMongo) CreatePlanet(planet domain.PlanetCreationObj, qtdFilms int) (*string, *errs.AppError) {
	var err error
	var ctx, _ = context.WithTimeout(context.Background(), 40*time.Second)

	// create the context and connection (once)
	client := Connect(r.Host)

	// select the db and collection
	collection := client.Conn.Database(r.Database).Collection("planets")

	insertResult, err := collection.InsertOne(ctx, planet)
	if err != nil {
		return nil, errs.NewUnexpectedError("Unexpected database error inserting planet")
	}
	oid, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errs.NewUnexpectedError("Unexpected database error planet id converting to mongo ObjectID")
	}

	stringOid := oid.Hex()
	return &stringOid, nil
}

func (r PlanetsRepositoryMongo) GetPlanetByName(name string) (*domain.Planet, *errs.AppError) {
	var err error
	var ctx, _ = context.WithTimeout(context.Background(), 40*time.Second)
	var planet domain.Planet

	// create the context and connection (once)
	client := Connect(r.Host)

	// select the db and collection
	collection := client.Conn.Database(r.Database).Collection("planets")

	filter := bson.D{
		{"name", name},
	}

	err = collection.FindOne(ctx, filter).Decode(&planet)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return nil, errs.NewNotFoundError("Planet with this name not found")
		default:
			return nil, errs.NewUnexpectedError("Unexpected database error find planet")
		}
	}

	return &planet, nil
}

func (r PlanetsRepositoryMongo) GetPlanetById(id string) (*domain.Planet, *errs.AppError) {
	var err error
	var ctx, _ = context.WithTimeout(context.Background(), 40*time.Second)
	var planet domain.Planet
	oid, err := convertStringToOID(id)
	if err != nil {
		logger.Error("Error converting id to mongo format")
		return nil, errs.NewBadRequestError("Error converting query parameter id to database mongo id format")
	}
	// create the context and connection (once)
	client := Connect(r.Host)

	// select the db and collection
	collection := client.Conn.Database(r.Database).Collection("planets")

	filter := bson.D{
		{"_id", oid},
	}

	err = collection.FindOne(ctx, filter).Decode(&planet)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return nil, errs.NewNotFoundError("Planet with this id not found")
		default:
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}
	}

	return &planet, nil
}

func (r PlanetsRepositoryMongo) DeletePlanetById(id string) *errs.AppError {
	var err error
	var ctx, _ = context.WithTimeout(context.Background(), 40*time.Second)
	oid, err := convertStringToOID(id)
	if err != nil {
		return errs.NewValidationError("Error converting query parameter id to database mongo id format") // no documents found
	}
	// create the context and connection (once)
	client := Connect(r.Host)

	// select the db and collection
	collection := client.Conn.Database(r.Database).Collection("planets")

	filter := bson.D{
		{"_id", oid},
	}

	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return errs.NewNotFoundError("Planet with this id not found")
		default:
			return errs.NewUnexpectedError("Unexpected database error")
		}
	}

	return nil
}

func convertStringToOID(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}
