package middleware

import (
	"github.com/KSkun/health-iot-backend/util"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			util.ResponseWrapper{Success: false, Message: "bad request", Error: err.Error()})
	}
	return nil
}

func NewValidator() *Validator {
	return &Validator{validator: validator.New()}
}
