package controller

import (
	"github.com/labstack/echo/v4"
	"log"
)

func InitController(echo *echo.Echo) {
	groupAPI := echo.Group("/api")

	groupAPIV1 := groupAPI.Group("/v1")
	initControllerV1(groupAPIV1)
	log.Printf("[Controller] Init done")
}

func initControllerV1(group *echo.Group) {
	groupUserV1 := group.Group("/user")
	initUserGroupV1(groupUserV1)
}
