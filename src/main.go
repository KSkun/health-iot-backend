package main

import (
	"fmt"
	"github.com/KSkun/health-iot-backend/config"
	"github.com/KSkun/health-iot-backend/controller"
	"github.com/KSkun/health-iot-backend/global"
	"github.com/KSkun/health-iot-backend/model"
	"github.com/KSkun/health-iot-backend/util"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
)

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			util.ResponseWrapper{Success: false, Message: "Bad Request", Error: err.Error()})
	}
	return nil
}

func main() {
	global.SetGlobalInit(echo.New())

	global.EchoInst.Validator = &Validator{validator: validator.New()}
	global.EchoInst.Use(echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
		Format:           "[Echo] ${time_custom} ${status} ${method} ${uri} ${latency_human} ${error}\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
	}))
	global.EchoInst.Use(echoMiddleware.Recover())
	global.EchoInst.Use(echoMiddleware.CORS())
	global.EchoInst.Debug = config.Debug

	config.InitConfig()
	model.InitMongo()
	controller.InitController(global.EchoInst)

	log.Printf("[Main] Starting server")
	addr := fmt.Sprintf("%s:%d", config.C.AppConfig.Addr, config.C.AppConfig.Port)
	err := global.EchoInst.Start(addr)
	if err != nil {
		log.Fatalf("[Main] Error when server running, %s", err.Error())
	}
}
