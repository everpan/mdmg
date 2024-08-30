package main

import (
	"github.com/everpan/mdmg/v8runtime"
	"github.com/everpan/mdmg/web/handler"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {
	defer v8runtime.DisposeCtxPool()

	app := fiber.New()
	// contrib/fiberzap
	logger, _ := zap.NewDevelopment()
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
	}))

	swgCfg := swagger.Config{FilePath: "./docs/swagger.json", Path: "./swagger"}
	app.Use(swagger.New(swgCfg))

	apiRouter := app.Group("/api")
	apiRouter.Group(handler.ICoderHandler.Path, handler.ICoderHandler.Handler)
	_ = app.Listen(":9091")
}
