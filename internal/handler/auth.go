package handler

import (
	"errors"
	"net/http"

	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/domain"
	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/service"
	"github.com/IvanMeln1k/go-online-trading-platform-app/pkg/tokens"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h *Handler) SignUp(ctx echo.Context) error {
	var body SignUpJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		logrus.Errorf("error bind user: %s", err)
		return echo.NewHTTPError(400, Message{Message: "Bad request"})
	}
	userId, err := h.services.Auth.SignUp(ctx.Request().Context(), domain.User{
		Email:    string(body.Email),
		Name:     body.Name,
		Username: body.Username,
		Password: body.Password,
	})
	if err != nil {
		logrus.Errorf("error create user: %s", err)
		if errors.Is(service.ErrEmailAlreadyInUse, err) {
			return echo.NewHTTPError(409, Message{Message: "Email already in use"})
		}
		if errors.Is(service.ErrUsernameAlreadyInUse, err) {
			return echo.NewHTTPError(409, Message{Message: "Username already in use"})
		}
		return echo.NewHTTPError(500, Message{Message: "Internal server error"})
	}
	return ctx.JSON(200, map[string]interface{}{"id": userId})
}

func (h *Handler) SignIn(ctx echo.Context) error {
	var body SignInJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(400, Message{Message: "Bad request"})
	}
	tokens, err := h.services.Auth.SignIn(ctx.Request().Context(), string(body.Email), body.Password)

	if err != nil {
		if errors.Is(err, service.ErrInvalidEmailOrPassowrd) {
			return echo.NewHTTPError(401, Message{Message: "Invalid username or password"})
		}
		return echo.NewHTTPError(401, Message{Message: "Unauthorized"})
	}

	ctx.SetCookie(&http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
	})

	return ctx.JSON(200, map[string]interface{}{"tokens": tokens})
}

func (h *Handler) Refresh(ctx echo.Context, params RefreshParams) error {
	tokens, err := h.services.Auth.Refresh(ctx.Request().Context(), params.RefreshToken.RefreshToken)

	if err != nil {
		if errors.Is(err, service.ErrSessionInvalidOrExpired) {
			return echo.NewHTTPError(401, Message{Message: "Unauthorized"})
		} else if errors.Is(err, service.ErrInternal) {
			return echo.NewHTTPError(500, Message{Message: "Internal server error"})
		}
		return echo.NewHTTPError(401, Message{Message: "Unauthorized"})
	}

	ctx.SetCookie(&http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
	})

	return ctx.JSON(200, map[string]interface{}{"tokens": tokens})
}

func (h *Handler) Logout(ctx echo.Context, params LogoutParams) error {
	err := h.services.Auth.Logout(ctx.Request().Context(), params.RefreshToken.RefreshToken)

	if err != nil {
		if errors.Is(err, service.ErrSessionInvalidOrExpired) {
			return echo.NewHTTPError(401, Message{Message: "Unauthorized"})
		}
		return echo.NewHTTPError(500, Message{Message: "Internal server error"})
	}

	return ctx.JSON(200, map[string]interface{}{"Message": "ok"})
}

func (h *Handler) LogoutAll(ctx echo.Context, params LogoutAllParams) error {
	err := h.services.Auth.LogoutAll(ctx.Request().Context(), params.RefreshToken.RefreshToken)

	if err != nil {
		if errors.Is(err, service.ErrSessionInvalidOrExpired) {
			return echo.NewHTTPError(401, Message{Message: "Unauthorized"})
		}
		return echo.NewHTTPError(500, Message{Message: "Internal server error"})
	}

	return ctx.JSON(200, map[string]interface{}{"Message": "ok"})
}

func (h *Handler) Verification(ctx echo.Context, params VerificationParams) error {
	emailToken := params.Token

	if *emailToken == "" {
		return echo.NewHTTPError(401, Message{Message: "No authorized"})
	}
	email, err := h.tokenManager.ParseEmailToken(*emailToken)
	if err != nil {
		if errors.Is(tokens.ErrTokenExpired, err) {
			return echo.NewHTTPError(401, Message{Message: "Token is expired"})
		}
		if errors.Is(tokens.ErrTokenInvalid, err) {
			return echo.NewHTTPError(401, Message{Message: "Token is invalid"})
		}
		return echo.NewHTTPError(500, Message{Message: "Internal server error"})
	}
	err = h.services.Auth.VerifyEmail(ctx.Request().Context(), email)
	if err != nil {
		if errors.Is(service.ErrUserNotFound, err) {
			return echo.NewHTTPError(401, Message{Message: "User not found"})
		}
		return echo.NewHTTPError(500, Message{Message: "Internal server error"})
	}
	return ctx.JSON(200, map[string]interface{}{"Message": "ok"})
}

func (h *Handler) ResendEmail(ctx echo.Context) error {
	id, err := h.getUserId(ctx)
	if err != nil {
		return err
	}

	err = h.services.Auth.ResendEmail(ctx.Request().Context(), id)
	if err != nil {
		logrus.Errorf("error send email verification: %s", err)
		return echo.NewHTTPError(500, Message{Message: "Internal server error"})
	}

	return ctx.JSON(200, map[string]interface{}{"Message": "ok"})
}

func (h *Handler) GetUser(ctx echo.Context) error {
	id, err := h.getUserId(ctx)
	if err != nil {
		return err
	}

	user, err := h.services.Auth.GetUser(ctx.Request().Context(), id)
	if err != nil {
		logrus.Errorf("get user error: %s", err)
		return echo.NewHTTPError(500, Message{Message: "Internal server error"})
	}

	return ctx.JSON(200, map[string]interface{}{"user": user})
}
