package controller

import "github.com/KSkun/health-iot-backend/model"

type ReqCreateDeviceV1 struct {
	Name   string `json:"name"`
	Serial string `json:"serial" validate:"required"`
}

type RspDeviceSimpleV1 struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Serial  string `json:"serial"`
	Online  bool   `json:"online"`
	Battery int    `json:"battery"`
	Warning bool   `json:"warning"`
}

func (r *RspDeviceSimpleV1) FromDeviceObject(o model.DeviceObject) {
	r.ID = o.ID.Hex()
	r.Name = o.Name
	r.Serial = o.Serial
	r.Online = o.IsOnline()
	r.Battery = o.Status.Battery
	r.Warning = len(o.Warnings()) > 0
}

type ReqTurnOffDeviceWarningV1 struct {
	IDHex string `param:"id" validate:"required"`
	Value int    `json:"value"`
}
