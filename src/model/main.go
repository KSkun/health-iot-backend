package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"time"
)

const defaultTimeout = 100 * time.Millisecond // 100 ms

type IModel interface {
	CreateUser(name string, password string) (primitive.ObjectID, error)
	GetUserByName(name string) (UserObject, bool, error)
	//CompareUserPassword(name string, password string) (bool, error)
}

var M IModel

func InitModel() {
	initMongo()
}

func defaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defaultTimeout)
}
