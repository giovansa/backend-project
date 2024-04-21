package handler

import (
	"github.com/SawitProRecruitment/UserService/internal"
	"github.com/SawitProRecruitment/UserService/model"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"net/http"
	"strings"
)

// (POST /hello)
func (s *Server) Register(ctx echo.Context) error {
	user := new(model.RegisterUserReq)
	if err := ctx.Bind(user); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	err := user.Validate()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	userInput, err := user.ToDAO()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	userID, err := s.Repository.RegisterUser(ctx.Request().Context(), userInput)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, map[string]interface{}{"data": map[string]string{
		"user_id": userID,
	}})
}

// (POST /login)
func (s *Server) Login(ctx echo.Context) error {
	req := new(model.LoginRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	ctx2 := ctx.Request().Context()
	userDAO, err := s.Repository.GetUserByPhone(ctx2, strings.TrimSpace(req.Phone))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	user := model.FromRepoUser(userDAO)
	err = user.CheckLogin(req.Password)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	tokenString, err := internal.GenerateJWTToken(user, s.Cfg.App.Env)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if err = s.Repository.IncrSuccessLogin(ctx2, user.Phone); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"token": tokenString,
	})
}

// (GET /profile)
func (s *Server) GetProfile(ctx echo.Context) error {
	claimUser := ctx.Get("claims").(*model.Claims)
	if claimUser == nil {
		return ctx.JSON(http.StatusForbidden, map[string]struct{}{})
	}

	userDAO, err := s.Repository.GetUserByPhone(ctx.Request().Context(), claimUser.Phone)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}
	user := model.FromRepoUser(userDAO)

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"data": user.ToProfileResp(),
	})
}

func (s *Server) UpdateUser(ctx echo.Context) error {
	claimUser := ctx.Get("claims").(*model.Claims)
	if claimUser == nil {
		return ctx.JSON(http.StatusForbidden, map[string]struct{}{})
	}

	updateUser := new(model.UpdateUserReq)
	if err := ctx.Bind(updateUser); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	err := updateUser.Validate()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	err = s.Repository.UpdateUser(ctx.Request().Context(), updateUser.ToDAO(), claimUser.Phone)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ctx.JSON(http.StatusConflict, map[string]string{"error": "invalid request"})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"data": updateUser,
	})
}
