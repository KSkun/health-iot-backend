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
	Field   string `json:"field"`
	Type    int    `json:"type"`
	Message string `json:"message"`
}

type DeviceStatusObject struct {
	Battery  int  `bson:"battery" json:"battery"`
	Locating bool `bson:"locating" json:"locating"`
	Wearing  bool `bson:"wearing" json:"wearing"`
}

type DeviceSensorObject struct {
	HeartRate   int     `bson:"heart_rate" json:"heart_rate"`
	BloodOxygen int     `bson:"blood_oxygen" json:"blood_oxygen"`
	Longitude   float32 `bson:"longitude" json:"longitude"`
	Latitude    float32 `bson:"latitude" json:"latitude"`
	SOSWarning  bool    `bson:"sos_warning" json:"sos_warning"`
	FallWarning bool    `bson:"fall_warning" json:"fall_warning"`
}

type DeviceObject struct {
	ID             primitive.ObjectID `bson:"_id" json:"-"` // Unique index
	IDHex          string             `bson:"-" json:"id"`
	Name           string             `bson:"name" json:"name"`
	Serial         string             `bson:"serial" json:"serial"` // Unique index
	OwnerID        primitive.ObjectID `bson:"owner_id" json:"-"`
	OwnerIDHex     string             `bson:"-" json:"owner_id"`
	LastReportTime int64              `bson:"last_report_time" json:"last_report_time"`
	Status         DeviceStatusObject `bson:"status" json:"status"`
	Sensor         DeviceSensorObject `bson:"sensor" json:"sensor"`
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
			DeviceWarning{Field: "blood_oxygen", Type: WarningTypeTooLow, Message: "blood oxygen saturation too low"})
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

func (o *DeviceObject) CompileJSON() {
	o.IDHex = o.ID.Hex()
	o.OwnerIDHex = o.OwnerID.Hex()
}
