package main

import (
	"fmt"
	"github.com/KSkun/health-iot-backend/config"
	"github.com/KSkun/health-iot-backend/global"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	global.SetGlobalInit(echo.New())

	global.EchoInst.Use(echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
		Format:           "[Echo] ${time_custom} ${status} ${method} ${uri} ${latency_human} ${error}\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
	}))
	global.EchoInst.Use(echoMiddleware.Recover())
	global.EchoInst.Use(echoMiddleware.CORS())
	global.EchoInst.Debug = config.Debug

	config.InitConfig()

	addr := fmt.Sprintf("%s:%d", config.C.AppConfig.Addr, config.C.AppConfig.Port)
	err := global.EchoInst.Start(addr)
	if err != nil {
		global.Logger.Fatalf("[Main] Error when server running, %s", err.Error())
	}
}
