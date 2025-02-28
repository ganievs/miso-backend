package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Services struct {
	Modules   string `json:"modules.v1"`
	Providers string `json:"providers.v1"`
}

func main() {
	e := echo.New()
	e.GET("/.well-known/terraform.json", func(c echo.Context) error {
		return c.JSON(http.StatusOK, &Services{
			Modules:   "/v1/modules/",
			Providers: "/v1/providers/",
		})
	})
	e.Logger.Fatal(e.Start(":1323"))
}
