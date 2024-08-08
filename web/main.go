package main

import (
	"github.com/everpan/mdmg/v8runtime"
	"github.com/everpan/mdmg/web/handler"
	"github.com/everpan/mdmg/web/zlog"
	"github.com/gofiber/fiber/v2"
)

func main() {
	zlog.InitLogger()
	defer zlog.Sync()
	defer v8runtime.DisposeCtxPool()

	app := fiber.New()
	apiRouter := app.Group("/api")
	apiRouter.Group(handler.ICoderHandler.Path, handler.ICoderHandler.Handler)
	app.Listen(":9091")
}
