package controller

import (
	"github.com/KSkun/health-iot-backend/controller/param"
	"github.com/KSkun/health-iot-backend/model"
	"github.com/KSkun/health-iot-backend/util"
	"github.com/labstack/echo/v4"
	"net/http"
)

func InitUserGroupV1(group *echo.Group) {
	group.POST("", HandlerCreateUserV1)
}

func HandlerCreateUserV1(ctx echo.Context) error {
	req := param.ReqCreateUserV1{}
	if err := ctx.Bind(&req); err != nil {
		return util.FailedResp(ctx, http.StatusBadRequest, "Bad Request", err.Error())
	}
	if err := ctx.Validate(req); err != nil {
		return util.FailedResp(ctx, http.StatusBadRequest, "Bad Request", err.Error())
	}

	// TODO Check if name already exists
	id, err := model.M.CreateUser(req.Name, req.Password)
	if err != nil {
		return util.FailedResp(ctx, http.StatusInternalServerError, "Database Error", err.Error())
	}
	return util.SuccessResp(ctx, http.StatusOK, echo.Map{"id": id})
}
