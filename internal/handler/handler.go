package handler

import (
	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/service"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *echo.Echo {
	router := echo.New()

	return router
}
