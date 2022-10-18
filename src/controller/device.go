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
