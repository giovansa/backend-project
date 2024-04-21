package transport

import (
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/labstack/echo/v4"
)

type API struct {
	Handler handler.HandlerInterface
	Router  *echo.Echo
}

func RegisterHandler(e *echo.Echo, handler handler.HandlerInterface) {
	e.POST("/user", handler.Register)
	e.PUT("/user", handler.UpdateUser, AuthMiddleware)
	e.POST("/login", handler.Login)
	e.GET("/profile", handler.GetProfile, AuthMiddleware)
}
