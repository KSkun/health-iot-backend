package model

import (
	"fmt"
	"github.com/KSkun/health-iot-backend/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

type mongoModel struct {
	client   *mongo.Client
	database *mongo.Database
	colUser  *mongo.Collection
}

func initMongo() {
	// Connect mongo
	ctx, cancel := defaultContext()
	defer cancel()
	uri := fmt.Sprintf("mongodb://%s:%d", config.C.MongoConfig.Addr, config.C.MongoConfig.Port)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("[Mongo] Error when connecting to mongo, %s", err.Error())
	}
	log.Printf("[Mongo] Connected to mongo %s", uri)
	// Check connectivity
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("[Mongo] Error when connecting to mongo, %s", err.Error())
	}
	log.Printf("[Mongo] Ping success to %s", uri)
	// Check if database exists
	exists, err := checkDatabaseExist(client)
	if err != nil {
		log.Fatalf("[Mongo] Error when checking mongo database, %s", err.Error())
	}
	// Init MongoModel
	db := client.Database(config.C.MongoConfig.Database)
	model := mongoModel{client: client, database: db}
	model.colUser = model.database.Collection("user")
	// Init mongo database
	if !exists {
		err = initMongoDatabase(&model)
		if err != nil {
			log.Fatalf("[Mongo] Error when initializing mongo database, %s", err.Error())
		}
	}

	M = &model
	log.Printf("[Model] Init done")
}

func checkDatabaseExist(client *mongo.Client) (bool, error) {
	// Check if database already exists, if false, do initialization
	ctx, cancel := defaultContext()
	defer cancel()
	dbNames, err := client.ListDatabaseNames(ctx, bson.M{}, &options.ListDatabasesOptions{})
	if err != nil {
		return false, err
	}
	found := false
	for i := 0; i < len(dbNames); i++ {
		if dbNames[i] == config.C.MongoConfig.Database {
			found = true
		}
	}
	return found, nil
}

func initMongoDatabase(m *mongoModel) error {
	ctx, cancel := defaultContext()
	defer cancel()
	_, err := m.colUser.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"name": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}
	return nil
}
