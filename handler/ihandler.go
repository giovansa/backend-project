package handler

import "github.com/labstack/echo/v4"

type HandlerInterface interface {
	Register(ctx echo.Context) error
	Login(ctx echo.Context) error
	GetProfile(ctx echo.Context) error
	UpdateUser(ctx echo.Context) error
}
