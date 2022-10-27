package model

import (
	"go.mongodb.org/mongo-driver/bson"
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
	GetDevice(id primitive.ObjectID) (DeviceObject, error)
	GetDeviceBySerial(serial string) (DeviceObject, bool, error)
	TurnOffDeviceWarning(id primitive.ObjectID) error
	AddReportData(deviceID primitive.ObjectID, time int64, status DeviceStatusObject, sensor bson.M) (primitive.ObjectID, error)
	GetReportDataByOwner(ownerID primitive.ObjectID, conditions bson.M) ([]ReportObject, error)
}

var M IModel

func InitModel() {
	initMongo()
}

func defaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defaultTimeout)
}
