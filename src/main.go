package main

import (
	"fmt"
	"github.com/KSkun/health-iot-backend/config"
	"github.com/KSkun/health-iot-backend/global"
	"github.com/KSkun/health-iot-backend/middleware"
	"github.com/KSkun/health-iot-backend/model"
	"github.com/KSkun/health-iot-backend/router"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"log"
)

func main() {
	global.SetGlobalInit(echo.New())

	global.EchoInst.Validator = middleware.NewValidator()
	global.EchoInst.Use(middleware.Logger())
	global.EchoInst.Use(echoMiddleware.Recover())
	global.EchoInst.Use(echoMiddleware.CORS())
	global.EchoInst.Debug = config.Debug

	config.InitConfig()
	model.InitModel()
	router.InitRouter(global.EchoInst)

	log.Printf("[Main] Starting server")
	addr := fmt.Sprintf("%s:%d", config.C.AppConfig.Addr, config.C.AppConfig.Port)
	err := global.EchoInst.Start(addr)
	if err != nil {
		log.Fatalf("[Main] Error when server running, %s", err.Error())
	}
}
