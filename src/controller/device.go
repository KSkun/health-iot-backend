package controller

import (
	"github.com/KSkun/health-iot-backend/controller/param"
	"github.com/KSkun/health-iot-backend/model"
	"github.com/KSkun/health-iot-backend/util"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func HandlerCreateDeviceV1(ctx echo.Context) error {
	req := param.ReqCreateDeviceV1{}
	if err := ctx.Bind(&req); err != nil {
		return util.FailedResp(ctx, http.StatusBadRequest, "bad request", err.Error())
	}
	if err := ctx.Validate(req); err != nil {
		return util.FailedResp(ctx, http.StatusBadRequest, "bad request", err.Error())
	}
	// Insert device to database
	userID := ctx.Get("id").(primitive.ObjectID)
	id, err := model.M.CreateDevice(req.Name, req.Serial, userID)
	if mongo.IsDuplicateKeyError(err) {
		return util.FailedResp(ctx, http.StatusBadRequest, "device already registered", err.Error())
	}
	if err != nil {
		return util.FailedResp(ctx, http.StatusInternalServerError, "database error", err.Error())
	}
	return util.SuccessResp(ctx, http.StatusOK, echo.Map{"id": id.Hex()})
}

func HandlerGetDevicesV1(ctx echo.Context) error {
	userID := ctx.Get("id").(primitive.ObjectID)
	devices_, err := model.M.GetDevicesByOwner(userID)
	if err != nil {
		return util.FailedResp(ctx, http.StatusInternalServerError, "database error", err.Error())
	}
	devices := []param.RspDeviceSimpleV1{}
	for _, d_ := range devices_ {
		d := param.RspDeviceSimpleV1{}
		d.FromDeviceObject(d_)
		devices = append(devices, d)
	}
	return util.SuccessResp(ctx, http.StatusOK, echo.Map{"devices": devices})
}

func HandlerGetDeviceV1(ctx echo.Context) error {
	// Convert hex string to ObjectId
	idHex := ctx.Param("id")
	if idHex == "" {
		return util.FailedResp(ctx, http.StatusBadRequest, "invalid device id", "empty id")
	}
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return util.FailedResp(ctx, http.StatusBadRequest, "invalid device id", err.Error())
	}
	// Query device and return
	userID := ctx.Get("id").(primitive.ObjectID)
	device, err := model.M.GetDevice(id)
	if err != nil {
		return util.FailedResp(ctx, http.StatusInternalServerError, "database error", err.Error())
	}
	if device.OwnerID != userID {
		return util.FailedResp(ctx, http.StatusForbidden, "forbidden", "device is not registered to user")
	}
	device.CompileJSON()
	return util.SuccessResp(ctx, http.StatusOK,
		echo.Map{"device": device, "warnings": device.Warnings(), "online": device.IsOnline()})
}
