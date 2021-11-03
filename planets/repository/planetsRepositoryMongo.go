package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
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

func (r PlanetsRepositoryMongo) GetAllPlanets() ([]domain.Planet, error) {
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
		return nil, err
	}

	if err := cur.All(ctx, &planets); err != nil {
		return nil, err
	}

	if planets == nil {
		planets = make([]domain.Planet, 0)
		return planets, err
	}

	return planets, err
}

func (r PlanetsRepositoryMongo) CreatePlanet(planetRequest domain.PlanetCreationRequest, qtdFilms int) (*domain.Planet, error) {
	var err error
	var ctx, _ = context.WithTimeout(context.Background(), 40*time.Second)

	// create the context and connection (once)
	client := Connect(r.Host)

	// select the db and collection
	collection := client.Conn.Database(r.Database).Collection("planets")

	planet := domain.PlanetCreationObj{
		Name:     planetRequest.Name,
		Climate:  planetRequest.Climate,
		Terrain:  planetRequest.Terrain,
		QtdFilms: qtdFilms,
	}

	insertResult, err := collection.InsertOne(ctx, planet)
	if err != nil {
		return nil, err
	}
	oid, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		err := errors.New("Error planet id converting to mongo ObjectID")
		return nil, err
	}

	insertedPlanet := domain.Planet{
		ID:       oid,
		Name:     planetRequest.Name,
		Climate:  planetRequest.Climate,
		Terrain:  planetRequest.Terrain,
		QtdFilms: qtdFilms,
	}

	return &insertedPlanet, err
}

func (r PlanetsRepositoryMongo) GetPlanetByName(name string) (*domain.Planet, error) {
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
			err := errors.New("Planet with this name not found")
			return nil, err // no documents found
		default:
			return nil, err // something else
		}
	}

	return &planet, err
}

func (r PlanetsRepositoryMongo) GetPlanetById(id string) (*domain.Planet, error) {
	var err error
	var ctx, _ = context.WithTimeout(context.Background(), 40*time.Second)
	var planet domain.Planet
	oid, err := convertStringToOID(id)
	if err != nil {
		err := errors.New("Error converting query parameter id to database id format")
		return nil, err // no documents found
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
			err := errors.New("Planet with this id not found")
			return nil, err // no documents found
		default:
			return nil, err // something else
		}
	}

	return &planet, err
}

func (r PlanetsRepositoryMongo) DeletePlanetById(id string) error {
	var err error
	var ctx, _ = context.WithTimeout(context.Background(), 40*time.Second)
	oid, err := convertStringToOID(id)
	if err != nil {
		err := errors.New("Error converting query parameter id to database id format")
		return err // no documents found
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
			err := errors.New("Planet with this id not found")
			return err // no documents found
		default:
			return err // something else
		}
	}

	return nil
}

func convertStringToOID(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}
