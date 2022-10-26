package router

import (
	"github.com/KSkun/health-iot-backend/controller"
	"github.com/KSkun/health-iot-backend/middleware"
	"github.com/labstack/echo/v4"
)

func initDeviceGroupV1(group *echo.Group) {
	group.POST("", controller.HandlerCreateDeviceV1, middleware.JWT)
	group.GET("/list", controller.HandlerGetDevicesV1, middleware.JWT)
	group.GET("/:id", controller.HandlerGetDeviceV1, middleware.JWT)
	group.PUT("/:id/warning", controller.HandlerTurnOffDeviceWarningV1, middleware.JWT)
	group.POST("/data", controller.HandlerAddReportDataV1) // for IoT device
}
