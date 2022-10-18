package controller

import (
	"github.com/KSkun/health-iot-backend/controller/param"
	"github.com/KSkun/health-iot-backend/model"
	"github.com/KSkun/health-iot-backend/util"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func HandlerCreateUserV1(ctx echo.Context) error {
	req := param.ReqUserSimpleV1{}
	if err := ctx.Bind(&req); err != nil {
		return util.FailedResp(ctx, http.StatusBadRequest, "bad request", err.Error())
	}
	if err := ctx.Validate(req); err != nil {
		return util.FailedResp(ctx, http.StatusBadRequest, "bad request", err.Error())
	}
	// Insert user to database
	id, err := model.M.CreateUser(req.Name, req.Password)
	if mongo.IsDuplicateKeyError(err) {
		return util.FailedResp(ctx, http.StatusBadRequest, "user name already exists", err.Error())
	}
	if err != nil {
		return util.FailedResp(ctx, http.StatusInternalServerError, "database error", err.Error())
	}
	return util.SuccessResp(ctx, http.StatusOK, echo.Map{"id": id.Hex()})
}

func HandlerLoginV1(ctx echo.Context) error {
	req := param.ReqUserSimpleV1{}
	if err := ctx.Bind(&req); err != nil {
		return util.FailedResp(ctx, http.StatusBadRequest, "bad request", err.Error())
	}
	if err := ctx.Validate(req); err != nil {
		return util.FailedResp(ctx, http.StatusBadRequest, "bad request", err.Error())
	}
	// Compare password
	user, found, err := model.M.GetUserByName(req.Name)
	if err != nil {
		return util.FailedResp(ctx, http.StatusInternalServerError, "database error", err.Error())
	}
	if !found {
		return util.FailedResp(ctx, http.StatusBadRequest, "user name and password do not match", "user not found")
	}
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return util.FailedResp(ctx, http.StatusBadRequest, "user name and password do not match", "wrong password")
	}
	if err != nil {
		return util.FailedResp(ctx, http.StatusInternalServerError, "internal server error", err.Error())
	}
	// Sign token
	tokenStr, expireTime, err := util.NewJWTToken(user.ID.Hex())
	if err != nil {
		return util.FailedResp(ctx, http.StatusInternalServerError, "internal server error", err.Error())
	}
	return util.SuccessResp(ctx, http.StatusOK, echo.Map{"token": tokenStr, "expire_time": expireTime.Unix()})
}
