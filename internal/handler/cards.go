package handler

import (
	"errors"
	"regexp"

	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/domain"
	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetAllCards(ctx echo.Context) error {
	id, err := h.getUserId(ctx)
	if err != nil {
		return err
	}
	cards, err := h.services.Cards.GetAll(ctx.Request().Context(), id)
	if err != nil {
		logrus.Errorf("error getting cards: %s", err)
		if errors.Is(service.ErrUserNotFound, err) {
			return echo.NewHTTPError(404, Message{Message: "User not found"})
		}
		return echo.NewHTTPError(500, Message{Message: "Internal server error"})
	}
	return ctx.JSON(200, map[string]interface{}{"cards": cards})
}

func (h *Handler) AddCard(ctx echo.Context) error {
	var CardData AddCardJSONRequestBody
	if err := ctx.Bind(&CardData); err != nil {
		return echo.NewHTTPError(500, Message{Message: "Internal server error"})
	}
	if re, err := regexp.Compile(`^\d{16}$`); !re.MatchString(CardData.Number) {
		if err != nil {
			return echo.NewHTTPError(500, Message{Message: "Internal server error"})
		}
		return echo.NewHTTPError(400, Message{Message: "Validation error"})
	}
	if re, err := regexp.Compile(`^\d{3}$`); !re.MatchString(CardData.Cvv) {
		if err != nil {
			return echo.NewHTTPError(500, Message{Message: "Internal server error"})
		}
		return echo.NewHTTPError(500, Message{Message: "Internal server error"})
	}
	if re, err := regexp.Compile(`([01][0-9]|2[0-9])\/[0-9]{2}$`); !re.MatchString(CardData.Data) {
		if err != nil {
			return echo.NewHTTPError(500, Message{Message: "Internal server error"})
		}
		return echo.NewHTTPError(400, Message{Message: "Validation error"})
	}
	id, err := h.getUserId(ctx)
	if err != nil {
		logrus.Errorf("Handler error: %s", err)
		return err
	}
	CardId, err := h.services.Cards.Create(ctx.Request().Context(), id, domain.Card{
		Id:     CardData.Id,
		Number: CardData.Number,
		Data:   CardData.Data,
		Cvv:    CardData.Cvv,
		UserId: CardData.UserId,
	})
	if err != nil {
		logrus.Errorf("Handler createCard error: %s", err)
		return echo.NewHTTPError(500, Message{Message: "Internal server error"})
	}
	return ctx.JSON(200, map[string]interface{}{"cardId": CardId})
}

func (h *Handler) DeleteCard(ctx echo.Context, cardId int) error {
	id, err := h.getUserId(ctx)
	if err != nil {
		logrus.Errorf("Handler deleteCard error: %s", err)
		return err
	}
	err = h.services.Cards.Delete(ctx.Request().Context(), id, cardId)
	if err != nil {
		logrus.Errorf("DeleteCard error in Handler: %s", err)
		return echo.NewHTTPError(500, Message{Message: "Internal server error"})
	}
	return ctx.JSON(200, map[string]interface{}{"status": "ok"})
}

func (h *Handler) GetTheCard(ctx echo.Context, cardId int) error {
	id, err := h.getUserId(ctx)
	if err != nil {
		logrus.Errorf("Handler getUserId error:%s", err)
		return err
	}
	card, err := h.services.Cards.Get(ctx.Request().Context(), id, cardId)
	if err != nil {
		if errors.Is(service.ErrCardNotFound, err) {
			logrus.Errorf("GetCard error in handler: %s", err)
			return newErrorResponse(404, "CardNotFound")
		}
		logrus.Errorf("GetCard error in handler: %s", err)
		return newErrorResponse(500, "Internal server error")
	}
	return ctx.JSON(200, card)
}
