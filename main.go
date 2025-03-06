package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"miso/internal/config"
)

type Services struct {
	Modules   string `json:"modules.v1"`
	Providers string `json:"providers.v1"`
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	config, err := config.LoadConfig()
	if err != nil {
		logger.Error("Could not load config")
	}

	// Main server
	mainServer := echo.New()
	mainServer.HideBanner = true
	mainServer.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))
	mainServer.Use(middleware.Recover())
	mainServer.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, true)
	})
	mainServer.GET("/.well-known/terraform.json", func(c echo.Context) error {
		return c.JSON(http.StatusOK, &Services{
			Modules:   "/v1/modules/",
			Providers: "/v1/providers/",
		})
	})

	// Health Rerver
	healthServer := echo.New()
	healthServer.HideBanner = true
	healthServer.Use(middleware.Logger())
	healthServer.Use(middleware.Recover())
	healthServer.Use(echoprometheus.NewMiddleware("miso"))

	healthServer.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, true)
	})
	healthServer.GET("/metrics", echoprometheus.NewHandler())

	// Run all
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := mainServer.Start(config.App.Host + ":" + config.App.Port); err != nil && err != http.ErrServerClosed {
			mainServer.Logger.Error("shutting down the server")
		}
	}()

	go func() {
		if err := healthServer.Start(config.App.Host + ":" + config.Metrics.Port); err != nil && err != http.ErrServerClosed {
			healthServer.Logger.Error("shutting down the server")
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := mainServer.Shutdown(shutdownCtx); err != nil {
		mainServer.Logger.Error(err)
	}

	if err := healthServer.Shutdown(shutdownCtx); err != nil {
		healthServer.Logger.Error(err)
	}
}
