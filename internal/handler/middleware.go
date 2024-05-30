package handler

import (
	"errors"
	"strings"

	"github.com/IvanMeln1k/go-online-trading-platform-app/pkg/tokens"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h *Handler) userIdentity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get("Authorization")

		if header == "" {
			logrus.Errorf("empty header auth")
			return newErrorResponse(401, "Not authorized")
		}

		params := strings.Split(header, " ")
		if len(params) < 2 {
			logrus.Errorf("invalid header auth: len(params) < 2")
			return newErrorResponse(401, "Not authorized")
		}
		if params[0] != "Bearer" {
			logrus.Errorf("invalid header auth: no bearer, but header: %s", header)
			return newErrorResponse(401, "Not authorized")
		}

		userId, role, err := h.tokenManager.ParseAccessToken(params[1])
		if err != nil {
			logrus.Errorf("error parsing token: %s", err)
			if errors.Is(tokens.ErrTokenExpired, err) {
				return newErrorResponse(401, "Token is expired")
			} else if errors.Is(tokens.ErrTokenInvalid, err) {
				return newErrorResponse(401, "Not authorized")
			}
			return newErrorResponse(500, "Internal server error")
		}

		c.Set("userId", userId)
		c.Set("role", role)

		return next(c)
	}
}

func (h *Handler) getUserId(c echo.Context) (int, error) {
	id := c.Get("userId")

	idInt, ok := id.(int)
	if !ok {
		return 0, newErrorResponse(401, "Unautharized")
	}

	return idInt, nil
}

func (h *Handler) getRole(c echo.Context) (string, error) {
	role := c.Get("role")

	roleString, ok := role.(string)
	if !ok {
		return "", newErrorResponse(401, "Unautharized")
	}

	return roleString, nil
}
