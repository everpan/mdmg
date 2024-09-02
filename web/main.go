package main

import (
	"github.com/everpan/mdmg/v8runtime"
	"github.com/everpan/mdmg/web/handler"
	"github.com/gofiber/contrib/fgprof"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"time"
)

func CreateApp() *fiber.App {
	app := fiber.New()
	// contrib/fiberzap
	logger, _ := zap.NewDevelopment()
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
	}))

	swgCfg := swagger.Config{FilePath: "./docs/swagger.json", Path: "./swagger"}
	app.Use(swagger.New(swgCfg))

	app.Use(fgprof.New())

	apiRouter := app.Group("/api")
	apiRouter.Group(handler.ICoderHandler.Path, handler.ICoderHandler.Handler)

	staticConf := fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		Index:         "index.html",
		CacheDuration: 10 * time.Second,
		MaxAge:        3600,
	}

	app.Static("/", "./public", staticConf)
	return app
}

func main() {
	defer v8runtime.DisposeCtxPool()

	app := CreateApp()
	app.Listen(":8080")
}
