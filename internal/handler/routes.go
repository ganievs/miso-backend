package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(v1 *echo.Group) {
	providers := v1.Group("/providers")
	providers.GET("/:namespace/:type/versions", h.ListProviderVersions)
	providers.GET("/:namespace/:type/:version/download/:os/:arch", h.DownloadProviderVersion)

	modules := v1.Group("/modules")
	modules.GET("/:namespace/:name/:provider/versions", h.ListModuleVersions)
	modules.GET("/:namespace/:name/:provider/:version/download", h.DownloadModuleVersion)

	mirror := v1.Group("/mirror")
	mirror.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, true)
	})
}
