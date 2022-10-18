package model

import (
	"fmt"
	"github.com/KSkun/health-iot-backend/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
)

const (
	colNameUser   = "user"
	colNameDevice = "device"
)

type mongoModel struct {
	client    *mongo.Client
	database  *mongo.Database
	colUser   *mongo.Collection
	colDevice *mongo.Collection
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
	model.colUser = model.database.Collection(colNameUser)
	model.colDevice = model.database.Collection(colNameDevice)
	// Init mongo database
	if !exists || os.Getenv("MONGO_FORCE_INIT") == "1" {
		log.Printf("[Mongo] Database does not exists, initializing")
		err = initMongoDatabase(&model)
		if err != nil {
			log.Fatalf("[Mongo] Error when initializing mongo database, %s", err.Error())
		}
		log.Printf("[Mongo] Init database done")
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
	log.Printf("[Mongo] Create unique index 'name' of collection 'user'")
	_, err := m.colUser.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"name": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}
	log.Printf("[Mongo] Create unique index 'serial' of collection 'device'")
	_, err = m.colDevice.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"serial": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}
	return nil
}
