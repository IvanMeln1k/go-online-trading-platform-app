package handler

import (
	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/service"
	"github.com/IvanMeln1k/go-online-trading-platform-app/pkg/validate"
	"github.com/go-playground/validator"
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

	router.Validator = &validate.CustomValidator{
		Validator: validator.New(),
	}

	router.GET("/", func(c echo.Context) error {
		return c.String(200, "hello")
	})

	return router
}
