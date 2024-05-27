package handler

import (
	"errors"

	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/domain"
	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetMyAllProducts(ctx echo.Context) error {
	id, err := h.getUserId(ctx)
	if err != nil {
		return err
	}
	products, err := h.services.Products.GetMyAll(ctx.Request().Context(), id)
	if err != nil {
		logrus.Errorf("error getting products: %s", err)
		if errors.Is(service.ErrUserNotFound, err) {
			return echo.NewHTTPError(404, Message{Message: "User not found"})
		}
		return echo.NewHTTPError(500, Message{Message: "Internal server error"})
	}
	return ctx.JSON(200, map[string]interface{}{"products": products})
}
func (h *Handler) GetAllProducts(ctx echo.Context, params GetAllProductsParams) error {
	products, err := h.services.Products.GetAll(ctx.Request().Context(), domain.Filter{
		Article:      *params.Params.Article,
		Name:         *params.Params.Name,
		MinPrice:     *params.Params.MinPrice,
		MaxPrice:     *params.Params.MaxPrice,
		Manufacturer: *params.Params.Manufacturer,
		Rating:       *params.Params.Rating,
		Limit:        *params.Params.Limit,
		Offset:       *params.Params.Offset,
	})
	if err != nil {
		logrus.Errorf("error getting products by filter: %s", err)
		return echo.NewHTTPError(500, Message{Message: "Internal server error"})
	}
	return ctx.JSON(200, map[string]interface{}{"products": products})
}

func (h *Handler) AddProduct(ctx echo.Context) error {
	var ProductData AddProductJSONRequestBody
	if err := ctx.Bind(&ProductData); err != nil {
		return echo.NewHTTPError(500, Message{Message: "Internal server error"})
	}

	id, err := h.getUserId(ctx)
	if err != nil {
		logrus.Errorf("Handler error: %s", err)
		return err
	}
	ProductId, err := h.services.Products.Create(ctx.Request().Context(), id, domain.Product{
		Id:           ProductData.Id,
		Article:      ProductData.Article,
		Name:         ProductData.Name,
		Price:        ProductData.Price,
		Manufacturer: ProductData.Manufacturer,
		SellerId:     ProductData.SellerId,
		Deleted:      ProductData.Deleted,
		Rating:       ProductData.Rating,
	})
	if err != nil {
		logrus.Errorf("Handler createProduct error: %s", err)
		return echo.NewHTTPError(500, Message{Message: "Internal server error"})
	}
	return ctx.JSON(200, map[string]interface{}{"productId": ProductId})
}

func (h *Handler) DeleteProduct(ctx echo.Context, productId int) error {
	id, err := h.getUserId(ctx)
	if err != nil {
		logrus.Errorf("Handler deleteProduct error: %s", err)
		return err
	}
	err = h.services.Products.Delete(ctx.Request().Context(), id, productId)
	if err != nil {
		logrus.Errorf("DeleteProduct error in Handler: %s", err)
		return echo.NewHTTPError(500, Message{Message: "Internal server error"})
	}
	return ctx.JSON(200, map[string]interface{}{"status": "ok"})
}

func (h *Handler) GetTheProduct(ctx echo.Context, productId int) error {
	id, err := h.getUserId(ctx)
	if err != nil {
		logrus.Errorf("Handler getUserId error:%s", err)
		return err
	}
	product, err := h.services.Products.Get(ctx.Request().Context(), id, productId)
	if err != nil {
		if errors.Is(service.ErrProductNotFound, err) {
			logrus.Errorf("GetProduct error in handler: %s", err)
			return newErrorResponse(404, "ProductNotFound")
		}
		logrus.Errorf("GetProduct error in handler: %s", err)
		return newErrorResponse(500, "Internal server error")
	}
	return ctx.JSON(200, product)
}
