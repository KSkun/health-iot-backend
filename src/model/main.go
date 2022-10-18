package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"time"
)

const defaultTimeout = 100 * time.Millisecond // 100 ms

type IModel interface {
	// user.go
	CreateUser(name string, password string) (primitive.ObjectID, error)
	GetUserByName(name string) (UserObject, bool, error)
	// device.go
	CreateDevice(name string, serial string, ownerID primitive.ObjectID) (primitive.ObjectID, error)
	GetDevicesByOwner(ownerID primitive.ObjectID) ([]DeviceObject, error)
}

var M IModel

func InitModel() {
	initMongo()
}

func defaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defaultTimeout)
}
