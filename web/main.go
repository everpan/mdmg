package main

import (
	"encoding/json"
	"fmt"
	"github.com/everpan/mdmg/utils"
	"github.com/everpan/mdmg/web/icode"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

var logger *zap.Logger

func main() {
	icode.InitLogger()
	defer logger.Sync()

	app := fiber.New()
	apiRouter := app.Group("/api")
	apiRouter.Group("/v1/icode/:module/:version?/:jsFile?/*", func(c *fiber.Ctx) error {
		fmt.Printf("%v | %v | %v | %v\n", c.Params("module"), c.Params("version"), c.Params("jsFile"), c.Queries())
		// c.Next()
		scriptFile := "web/handler/test_001.mjs"
		// zCtx := zcode.GetCtx(scriptFile, c)
		zCtx := icode.GetCtx(c)
		script, e := os.ReadFile(scriptFile)
		if e != nil {
			return c.SendString(e.Error())
		}
		scriptFile = filepath.Base(scriptFile)
		r, e := zCtx.V8Ctx.RunScript(string(script), scriptFile)
		if e != nil {
			return c.SendString(e.Error())
		}
		v, e := utils.ToGoValue(zCtx.V8Ctx, r)
		if e != nil {
			return c.SendString(e.Error())
		}
		// fmt.Printf("go value:%v", v)
		jv, e := json.Marshal(v)
		if e != nil {
			return c.SendString(e.Error())
		}
		return c.Send(jv)
	})
	app.Listen(":9091")
}
