package handler

import (
	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/service"
	"github.com/IvanMeln1k/go-online-trading-platform-app/pkg/tokens"
	"github.com/IvanMeln1k/go-online-trading-platform-app/pkg/validate"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	services     *service.Service
	tokenManager tokens.TokenManagerI
}

type Deps struct {
	Services     *service.Service
	TokenManager tokens.TokenManagerI
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		services:     deps.Services,
		tokenManager: deps.TokenManager,
	}
}

func (h *Handler) InitRoutes() *echo.Echo {
	router := echo.New()

	router.Validator = &validate.CustomValidator{
		Validator: validator.New(),
	}

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/refresh", h.refresh)
		auth.DELETE("/logout", h.logout)
		auth.DELETE("/logout-all", h.logoutAll)
		auth.POST("/verify-email", h.verifyEmail)
		auth.POST("/resend-email", h.resendEmail, h.userIdentity)
	}

	return router
}
