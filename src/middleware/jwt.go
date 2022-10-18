package middleware

import (
	"github.com/KSkun/health-iot-backend/util"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
)

func JWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		tokenStr := ctx.Request().Header.Get("Authorization")
		if !strings.HasPrefix(tokenStr, "Bearer ") {
			return util.FailedResp(ctx, http.StatusUnauthorized, "bad token", "token not set")
		}
		tokenStr = strings.Replace(tokenStr, "Bearer ", "", 1)
		idHex, err := util.ValidateJWTToken(tokenStr)
		if err != nil {
			return util.FailedResp(ctx, http.StatusUnauthorized, "bad token", err.Error())
		}
		id, err := primitive.ObjectIDFromHex(idHex)
		if err != nil {
			return util.FailedResp(ctx, http.StatusUnauthorized, "bad token", err.Error())
		}
		ctx.Set("id", id)
		return next(ctx)
	}
}
