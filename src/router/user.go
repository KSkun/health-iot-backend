package router

import (
	"github.com/KSkun/health-iot-backend/controller"
	"github.com/labstack/echo/v4"
)

func initUserGroupV1(group *echo.Group) {
	group.POST("", controller.HandlerCreateUserV1)
	group.GET("/token", controller.HandlerLoginV1)
}
