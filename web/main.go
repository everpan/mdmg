package main

import (
	"github.com/everpan/mdmg/web/handler"
	"github.com/everpan/mdmg/web/icode/v8runtime"
	"github.com/everpan/mdmg/web/web_logger"
	"github.com/gofiber/fiber/v2"
)

func main() {
	web_logger.InitLogger()
	defer web_logger.Sync()
	defer v8runtime.DisposeCtxPool()

	app := fiber.New()
	apiRouter := app.Group("/api")
	apiRouter.Group(handler.ICoderHandler.Path, handler.ICoderHandler.Handler)
	app.Listen(":9091")
}
