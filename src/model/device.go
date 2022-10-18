package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeviceStatusObject struct {
	Battery  float32 `bson:"battery"`
	Locating bool    `bson:"locating"`
	Wearing  bool    `bson:"wearing"`
}

type DeviceSensorObject struct {
	HeartRate   int     `bson:"heart_rate"`
	BloodOxygen int     `bson:"blood_oxygen"`
	Longitude   float32 `bson:"longitude"`
	Latitude    float32 `bson:"latitude"`
	SOSWarning  bool    `bson:"sos_warning"`
	FallWarning bool    `bson:"fall_warning"`
}

type DeviceObject struct {
	ID             primitive.ObjectID `bson:"_id"` // Unique index
	Name           string             `bson:"name"`
	Serial         string             `bson:"serial"` // Unique index
	OwnerID        primitive.ObjectID `bson:"owner_id"`
	LastReportTime int64              `bson:"last_report_time"`
	Status         DeviceStatusObject `bson:"status"`
	Sensor         DeviceSensorObject `bson:"sensor"`
}

type ReportObject struct {
	ID       primitive.ObjectID `bson:"_id"`
	DeviceID primitive.ObjectID `bson:"device_id"`
	Time     int64              `bson:"time"`
	Status   DeviceStatusObject `bson:"status"`
	Sensor   bson.M             `bson:"sensor"`
}

func (m *mongoModel) CreateDevice(name string, serial string, ownerID primitive.ObjectID) (primitive.ObjectID, error) {
	ctx, cancel := defaultContext()
	defer cancel()
	res, err := m.colDevice.InsertOne(ctx,
		DeviceObject{ID: primitive.NewObjectID(), Name: name, Serial: serial, OwnerID: ownerID})
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}
