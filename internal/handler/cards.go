package handler

import (
	"errors"

	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/domain"
	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/service"
	"github.com/go-delve/delve/service"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type GetAllCardsReturn struct {
	cards []domain.CardReturn
}

func (h *Handler) getAllCards(c echo.Context) error {
	id, err := h.getUserId(c)
	if err != nil {
		return err
	}
	cards, err := h.services.Users.GetAllCards(c.Request().Context(), id)
	if err != nil {
		logrus.Errorf("error getting cards: %s", err)
		if errors.Is(service.ErrUserNotFound, err) {
			return newErrorResponse(404, "User not found")
		}
		return newErrorResponse(500, "Internal server error")
	}
	return c.JSON(200, GetAllCardsReturn{
		cards: cards,
	})

}

func (h *Handler) getCard(c echo.Context) error {
	return c.JSON(200, nil)
}

func (h *Handler) addCard(c echo.Context) error {
	return c.JSON(200, nil)
}

func (h *Handler) deleteCard(c echo.Context) error {
	return c.JSON(200, nil)
}
