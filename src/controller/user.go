package controller

import (
	"github.com/KSkun/health-iot-backend/controller/param"
	"github.com/KSkun/health-iot-backend/model"
	"github.com/KSkun/health-iot-backend/util"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func initUserGroupV1(group *echo.Group) {
	group.POST("", handlerCreateUserV1)
}

func handlerCreateUserV1(ctx echo.Context) error {
	req := param.ReqCreateUserV1{}
	if err := ctx.Bind(&req); err != nil {
		return util.FailedResp(ctx, http.StatusBadRequest, "Bad Request", err.Error())
	}
	if err := ctx.Validate(req); err != nil {
		return util.FailedResp(ctx, http.StatusBadRequest, "Bad Request", err.Error())
	}

	// Insert user to database
	id, err := model.M.CreateUser(req.Name, req.Password)
	if mongo.IsDuplicateKeyError(err) {
		return util.FailedResp(ctx, http.StatusBadRequest, "User name already exists", err.Error())
	}
	if err != nil {
		return util.FailedResp(ctx, http.StatusInternalServerError, "Database Error", err.Error())
	}
	return util.SuccessResp(ctx, http.StatusOK, echo.Map{"id": id.Hex()})
}
