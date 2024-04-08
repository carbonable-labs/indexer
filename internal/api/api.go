package api

import (
	"context"
	"net/http"

	"github.com/carbonable-labs/indexer/internal/config"
	"github.com/carbonable-labs/indexer/internal/storage"
	"github.com/labstack/echo/v4"
)

func Run(ctx context.Context, storage storage.Storage) {
	e := echo.New()
	cr := config.NewPebbleContractRepository(storage)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/register", func(c echo.Context) error {
		config := config.NewCongig()
		err := c.Bind(config)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}
		config = config.ComputeHash()

		err = cr.SaveConfig(*config)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"app_name": config.AppName,
			"hash":     config.Hash,
		})
	})

	e.Logger.Fatal(e.Start(":8080"))

	<-ctx.Done()
	_ = e.Shutdown(ctx)
}
