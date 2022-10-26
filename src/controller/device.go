package controller

import (
	"github.com/KSkun/health-iot-backend/model"
	"github.com/KSkun/health-iot-backend/util"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

func HandlerCreateDeviceV1(ctx echo.Context) error {
	req := ReqCreateDeviceV1{}
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
	devices := []RspDeviceSimpleV1{}
	for _, d_ := range devices_ {
		d := RspDeviceSimpleV1{}
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

func HandlerTurnOffDeviceWarningV1(ctx echo.Context) error {
	req := ReqTurnOffDeviceWarningV1{}
	if err := ctx.Bind(&req); err != nil {
		return util.FailedResp(ctx, http.StatusBadRequest, "bad request", err.Error())
	}
	if err := ctx.Validate(req); err != nil {
		return util.FailedResp(ctx, http.StatusBadRequest, "bad request", err.Error())
	}
	// Turning off needs a `false` value
	if req.Value != 0 {
		return util.FailedResp(ctx, http.StatusNotImplemented, "request parameter unsupported", "value param must be 0")
	}
	// Convert hex string to ObjectId
	id, err := primitive.ObjectIDFromHex(req.IDHex)
	if err != nil {
		return util.FailedResp(ctx, http.StatusBadRequest, "invalid device id", err.Error())
	}
	// Update device warning status
	userID := ctx.Get("id").(primitive.ObjectID)
	device, err := model.M.GetDevice(id)
	if err != nil {
		return util.FailedResp(ctx, http.StatusInternalServerError, "database error", err.Error())
	}
	if device.OwnerID != userID {
		return util.FailedResp(ctx, http.StatusForbidden, "forbidden", "device is not registered to user")
	}
	err = model.M.TurnOffDeviceWarning(id)
	if err != nil {
		return util.FailedResp(ctx, http.StatusInternalServerError, "database error", err.Error())
	}
	return util.SuccessResp(ctx, http.StatusOK, true)
}

func HandlerAddReportDataV1(ctx echo.Context) error {
	req := ReqAddReportDataV1{}
	if err := ctx.Bind(&req); err != nil {
		return util.FailedResp(ctx, http.StatusBadRequest, "bad request", err.Error())
	}
	if err := ctx.Validate(req); err != nil {
		return util.FailedResp(ctx, http.StatusBadRequest, "bad request", err.Error())
	}
	// Query device ID
	device, found, err := model.M.GetDeviceBySerial(req.Serial)
	if err != nil {
		return util.FailedResp(ctx, http.StatusInternalServerError, "database error", err.Error())
	}
	if !found {
		return util.FailedResp(ctx, http.StatusForbidden, "device not registered", "")
	}
	// Insert report data
	reportID, err := model.M.AddReportData(device.ID, time.Now().UnixMilli(), req.Status, bson.M(req.Sensor))
	if err != nil {
		return util.FailedResp(ctx, http.StatusInternalServerError, "database error", err.Error())
	}
	return util.SuccessResp(ctx, http.StatusOK, echo.Map{"id": reportID.Hex()})
}
