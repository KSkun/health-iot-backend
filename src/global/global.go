package global

import (
	"github.com/labstack/echo/v4"
)

var EchoInst *echo.Echo

// SetGlobalInit Init global variables
// **Call this before doing any init including config**
func SetGlobalInit(echo_ *echo.Echo) {
	EchoInst = echo_
}
