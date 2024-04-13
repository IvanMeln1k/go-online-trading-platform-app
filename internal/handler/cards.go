package handler

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/domain"
	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type GetAllCardsReturn struct {
	cards []domain.Card
}

func (h *Handler) getAllCards(c echo.Context) error {
	id, err := h.getUserId(c)
	if err != nil {
		return err
	}
	cards, err := h.services.Cards.GetAll(c.Request().Context(), id)
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
	cardIdstr := c.Param("cardId")
	cardId, err := strconv.Atoi(cardIdstr)
	if err != nil {
		logrus.Errorf("Convert string to int error in handler: %s", err)
		return newErrorResponse(500, "Internal server error")
	}
	id, err := h.getUserId(c)
	if err != nil {
		logrus.Errorf("Handler getUserId error:%s", err)
		return err
	}
	card, err := h.services.Cards.Get(c.Request().Context(), id, cardId)
	if err != nil {
		if errors.Is(service.ErrCardNotFound, err) {
			logrus.Errorf("GetCard error in handler: %s", err)
			return newErrorResponse(404, "CardNotFound")
		}
		logrus.Errorf("GetCard error in handler: %s", err)
		return newErrorResponse(500, "Internal server error")
	}
	return c.JSON(200, card)
}

func (h *Handler) addCard(c echo.Context) error {
	var CardData domain.Card
	if err := c.Bind(&CardData); err != nil {
		return newErrorResponse(500, "Internal server error")
	}
	if re, err := regexp.Compile(`^\d{16}$`); !re.MatchString(CardData.Number) {
		if err != nil {
			return newErrorResponse(500, "Validation error")
		}
		return newErrorResponse(400, "Validation error")
	}
	if re, err := regexp.Compile(`^\d{3}$`); !re.MatchString(CardData.Cvv) {
		if err != nil {
			return newErrorResponse(500, "Validation error")
		}
		return newErrorResponse(400, "Validation error")
	}
	if re, err := regexp.Compile(`([01][0-9]|2[0-9])\/[0-9]{2}$`); !re.MatchString(CardData.Data) {
		if err != nil {
			return newErrorResponse(500, "Validation error")
		}
		return newErrorResponse(400, "Validation error")
	}
	id, err := h.getUserId(c)
	if err != nil {
		logrus.Errorf("Handler error: %s", err)
		return err
	}
	CardId, err := h.services.Cards.Create(c.Request().Context(), id, CardData)
	if err != nil {
		logrus.Errorf("Handler createCard error: %s", err)
		return newErrorResponse(500, "Internal server error")
	}
	return c.JSON(200, CardId)
}

func (h *Handler) deleteCard(c echo.Context) error {
	cardIdstr := c.Param("cardId")
	cardId, err := strconv.Atoi(cardIdstr)
	if err != nil {
		logrus.Errorf("Convert string to int error in handler: %s", err)
		return newErrorResponse(500, "Internal server error")
	}
	id, err := h.getUserId(c)
	if err != nil {
		logrus.Errorf("Handler deleteCard error: %s", err)
		return err
	}
	err = h.services.Cards.Delete(c.Request().Context(), id, cardId)
	if err != nil {
		logrus.Errorf("DeleteCard error in Handler: %s", err)
		return newErrorResponse(500, "Internal server error")
	}
	return c.JSON(200, echo.Map{
		"condition": "success",
	})
}
