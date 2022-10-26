package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func (m *mongoModel) GetDeviceBySerial(serial string) (DeviceObject, bool, error) {
	ctx, cancel := defaultContext()
	defer cancel()
	res := m.colDevice.FindOne(ctx, bson.M{"serial": serial})
	if res.Err() == mongo.ErrNoDocuments {
		return DeviceObject{}, false, nil
	}
	if res.Err() != nil {
		return DeviceObject{}, false, res.Err()
	}
	device := DeviceObject{}
	err := res.Decode(&device)
	if err != nil {
		return DeviceObject{}, false, err
	}
	return device, true, nil
}

func (m *mongoModel) TurnOffDeviceWarning(id primitive.ObjectID) error {
	ctx, cancel := defaultContext()
	defer cancel()
	_, err := m.colDevice.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"warning": false}})
	return err
}

func (m *mongoModel) AddReportData(deviceID primitive.ObjectID, time int64, status DeviceStatusObject, sensor bson.M) (primitive.ObjectID, error) {
	ctx, cancel := defaultContext()
	defer cancel()
	reportID := primitive.ObjectID{}
	err := m.client.UseSession(ctx, func(sessionCtx mongo.SessionContext) error {
		err := sessionCtx.StartTransaction()
		if err != nil {
			return err
		}
		// Insert report object
		res, err := m.colReport.InsertOne(sessionCtx, ReportObject{
			ID:       primitive.NewObjectID(),
			DeviceID: deviceID,
			Time:     time,
			Status:   status,
			Sensor:   sensor,
		})
		if err != nil {
			sessionCtx.AbortTransaction(sessionCtx)
			return err
		}
		reportID = res.InsertedID.(primitive.ObjectID)
		// Update device object
		updateFields := bson.M{"status": status}
		for k, v := range sensor {
			updateFields["sensor."+k] = v
		}
		_, err = m.colDevice.UpdateOne(ctx, bson.M{"_id": deviceID}, bson.M{"$set": updateFields})
		if err != nil {
			sessionCtx.AbortTransaction(sessionCtx)
			return err
		}
		err = sessionCtx.CommitTransaction(sessionCtx)
		if err != nil {
			return err
		}
		return nil
	})
	return reportID, err
}
