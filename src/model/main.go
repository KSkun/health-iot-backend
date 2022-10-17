package model

import (
	"fmt"
	"github.com/KSkun/health-iot-backend/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/net/context"
	"log"
	"time"
)

const defaultTimeout = 100 * time.Millisecond // 100 ms

type MongoModel struct {
	client   *mongo.Client
	database *mongo.Database
	colUser  *mongo.Collection
}

type IModel interface {
	CreateUser(name string, password string) (string, error)
	//GetUser(id string) (UserObject, error)
	//GetUserByName(name string) (UserObject, bool, error)
	//CompareUserPassword(name string, password string) (bool, error)
}

var M IModel

func InitMongo() {
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
	// Init MongoModel
	db := client.Database(config.C.MongoConfig.Database)
	model := MongoModel{client: client, database: db}
	model.colUser = model.database.Collection("user")
	M = &model
	log.Printf("[Model] Init done")
}

func defaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defaultTimeout)
}
