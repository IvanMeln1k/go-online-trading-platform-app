package handler

import (
	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/service"
	"github.com/IvanMeln1k/go-online-trading-platform-app/pkg/tokens"
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

	RegisterHandlers(router, h)

	return router
}
