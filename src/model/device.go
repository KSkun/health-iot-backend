package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

func (m *mongoModel) GetDevicesByOwner(ownerID primitive.ObjectID) ([]DeviceObject, error) {
	ctx, cancel := defaultContext()
	defer cancel()
	res, err := m.colDevice.Find(ctx, bson.M{"owner_id": ownerID})
	if err != nil {
		return nil, err
	}
	devices := []DeviceObject{}
	err = res.All(ctx, &devices)
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (m *mongoModel) GetDevice(id primitive.ObjectID) (DeviceObject, error) {
	ctx, cancel := defaultContext()
	defer cancel()
	res := m.colDevice.FindOne(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		return DeviceObject{}, res.Err()
	}
	device := DeviceObject{}
	err := res.Decode(&device)
	if err != nil {
		return DeviceObject{}, err
	}
	return device, nil
}

func (m *mongoModel) TurnOffDeviceWarning(id primitive.ObjectID) error {
	ctx, cancel := defaultContext()
	defer cancel()
	_, err := m.colDevice.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"warning": false}})
	return err
}
