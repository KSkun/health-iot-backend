package util

import (
	"github.com/KSkun/health-iot-backend/config"
	"github.com/labstack/echo/v4"
)

type ResponseWrapper struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Error   string      `json:"error"`
}

func SuccessResp(ctx echo.Context, code int, object interface{}) error {
	return ctx.JSON(code, ResponseWrapper{Success: true, Data: object})
}

func FailedResp(ctx echo.Context, code int, message string, error string) error {
	if !config.Debug {
		error = ""
	}
	return ctx.JSON(code, ResponseWrapper{Success: false, Message: message, Error: error})
}
