package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(v1 *echo.Group) {
	providers := v1.Group("/providers")
	providers.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, true)
	})

	modules := v1.Group("/providers")
	modules.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, true)
	})

	mirror := v1.Group("/mirror")
	mirror.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, true)
	})

}
