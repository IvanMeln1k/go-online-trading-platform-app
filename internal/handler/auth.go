package handler

import (
	"errors"
	"net/http"

	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/domain"
	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type SignUpReturn struct {
	Id int `json:"id"`
}

func (h *Handler) signUp(c echo.Context) error {
	var user domain.User

	if err := c.Bind(&user); err != nil {
		logrus.Errorf("error bind user: %s", err)
		return newErrorResponse(400, "Bad request")
	}
	if err := c.Validate(user); err != nil {
		logrus.Errorf("error validate user: %s", err)
		return newErrorResponse(400, err.Error())
	}

	userId, err := h.services.Auth.SignUp(c.Request().Context(), user)
	if err != nil {
		logrus.Errorf("error create user: %s", err)
		if errors.Is(service.ErrEmailAlreadyInUse, err) {
			return newErrorResponse(409, "Email already in use")
		}
		if errors.Is(service.ErrUsernameAlreadyInUse, err) {
			return newErrorResponse(409, "Username already in use")
		}
		return newErrorResponse(500, "Internal server error")
	}

	return c.JSON(200, SignUpReturn{
		Id: userId,
	})
}

type signInInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type signInReturn struct {
	Tokens domain.Tokens `json:"tokens"`
}

func (h *Handler) signIn(c echo.Context) error {
	user := new(signInInput)
	if err := c.Bind(user); err != nil {
		return newErrorResponse(400, err.Error())
	}
	if err := c.Validate(user); err != nil {
		return newErrorResponse(400, err.Error())
	}

	tokens, err := h.services.Auth.SignIn(c.Request().Context(), user.Email, user.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidEmailOrPassowrd) {
			return newErrorResponse(401, "Invalid username or password")
		}
		return newErrorResponse(401, "Anauthorized")
	}

	c.SetCookie(&http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
	})
	return c.JSON(200, signInReturn{
		Tokens: tokens,
	})
}

type RefreshReturn struct {
	Tokens domain.Tokens `json:"tokens"`
}

func (h *Handler) refresh(c echo.Context) error {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		return newErrorResponse(401, "Unauthorized")
	}
	tokens, err := h.services.Auth.Refresh(c.Request().Context(), refreshToken.Value)
	if err != nil {
		if errors.Is(err, service.ErrSessionInvalidOrExpired) {
			return newErrorResponse(401, "Unauthorized")
		} else if errors.Is(err, service.ErrInternal) {
			return newErrorResponse(500, "Internal server error")
		}
		return newErrorResponse(401, "Unauthorized")
	}
	c.SetCookie(&http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
	})
	return c.JSON(200, RefreshReturn{
		Tokens: tokens,
	})
}

type LogoutReturn struct {
	Status string `json:"status"`
}

func (h *Handler) logout(c echo.Context) error {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		return newErrorResponse(401, "Unauthorized")
	}
	err = h.services.Auth.Logout(c.Request().Context(), refreshToken.Value)
	if err != nil {
		if errors.Is(err, service.ErrSessionInvalidOrExpired) {
			return newErrorResponse(401, "Unauthorized")
		}
		return newErrorResponse(500, "Internal server error")
	}
	return c.JSON(200, LogoutReturn{
		Status: "ok",
	})
}

func (h *Handler) logoutAll(c echo.Context) error {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		return newErrorResponse(401, "Unauthorized")
	}
	err = h.services.Auth.LogoutAll(c.Request().Context(), refreshToken.Value)
	if err != nil {
		if errors.Is(err, service.ErrSessionInvalidOrExpired) {
			return newErrorResponse(401, "Unauthorized")
		}
		return newErrorResponse(500, "Internal server error")
	}
	return c.JSON(200, LogoutReturn{
		Status: "ok",
	})
}
