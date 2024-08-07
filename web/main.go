package main

import (
	"github.com/everpan/mdmg/web/icode"
	"github.com/everpan/mdmg/web/icode/v8runtime"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var logger *zap.Logger

func main() {
	icode.InitLogger()
	defer logger.Sync()
	defer v8runtime.DisposeCtxPool()

	app := fiber.New()
	apiRouter := app.Group("/api")
	apiRouter.Group("/v1/icode/:modVer/:jsFile/*", icodeHandler)
	app.Listen(":9091")
}
