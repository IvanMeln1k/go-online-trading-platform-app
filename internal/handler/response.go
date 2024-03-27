package handler

import "github.com/labstack/echo/v4"

func newErrorResponse(statusCode int, message string) error {
	return echo.NewHTTPError(statusCode, message)
}
