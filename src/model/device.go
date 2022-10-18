package model

import (
	"github.com/KSkun/health-iot-backend/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	WarningTypeWarning = 1
	WarningTypeTooLow  = 2
	WarningTypeTooHigh = 3
)

type DeviceWarning struct {
	Field   string
	Type    int
	Message string
}

type DeviceStatusObject struct {
	Battery  int  `bson:"battery"`
	Locating bool `bson:"locating"`
	Wearing  bool `bson:"wearing"`
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

func (o DeviceStatusObject) Warnings() []DeviceWarning {
	warnings := []DeviceWarning{}
	if o.Battery < config.C.AppConfig.BatteryThreshold {
		warnings = append(warnings,
			DeviceWarning{Field: "battery", Type: WarningTypeTooLow, Message: "low battery"})
	}
	return warnings
}

func (o DeviceSensorObject) Warnings() []DeviceWarning {
	warnings := []DeviceWarning{}
	if o.HeartRate < config.C.AppConfig.HeartRateThreshold.Low {
		warnings = append(warnings,
			DeviceWarning{Field: "heart_rate", Type: WarningTypeTooLow, Message: "heart rate too low"})
	}
	if o.HeartRate > config.C.AppConfig.HeartRateThreshold.High {
		warnings = append(warnings,
			DeviceWarning{Field: "heart_rate", Type: WarningTypeTooHigh, Message: "heart rate too high"})
	}
	if o.BloodOxygen < config.C.AppConfig.BloodOxygenThreshold {
		warnings = append(warnings,
			DeviceWarning{Field: "blood_oxygen", Type: WarningTypeTooHigh, Message: "blood oxygen saturation too low"})
	}
	if o.FallWarning {
		warnings = append(warnings,
			DeviceWarning{Field: "fall_warning", Type: WarningTypeWarning, Message: "fall warning"})
	}
	if o.SOSWarning {
		warnings = append(warnings,
			DeviceWarning{Field: "sos_warning", Type: WarningTypeWarning, Message: "sos"})
	}
	return warnings
}

func (o DeviceObject) Warnings() []DeviceWarning {
	warnings := []DeviceWarning{}
	warnings = append(warnings, o.Status.Warnings()...)
	warnings = append(warnings, o.Sensor.Warnings()...)
	return warnings
}

func (o DeviceObject) IsOnline() bool {
	return time.Now().Before(time.Unix(o.LastReportTime, 0).Add(config.C.AppConfig.OnlineTimeoutDuration))
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
