package router

import (
	"github.com/labstack/echo/v4"
	"log"
)

func InitRouter(echo *echo.Echo) {
	groupAPI := echo.Group("/api")

	groupAPIV1 := groupAPI.Group("/v1")
	initAPIGroupV1(groupAPIV1)
	log.Printf("[Router] Init done")
}

func initAPIGroupV1(group *echo.Group) {
	groupUserV1 := group.Group("/user")
	initUserGroupV1(groupUserV1)
	groupDeviceV1 := group.Group("/device")
	initDeviceGroupV1(groupDeviceV1)
}
