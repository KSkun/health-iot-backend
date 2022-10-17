package controller

import (
	"github.com/labstack/echo/v4"
	"log"
)

func InitController(echo *echo.Echo) {
	groupAPI := echo.Group("/api")

	groupAPIV1 := groupAPI.Group("/v1")
	InitControllerV1(groupAPIV1)
	log.Printf("[Controller] Init done")
}

func InitControllerV1(group *echo.Group) {
	groupUserV1 := group.Group("/user")
	InitUserGroupV1(groupUserV1)
}
