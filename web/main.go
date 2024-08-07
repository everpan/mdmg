package main

import (
	"encoding/json"
	"fmt"
	"github.com/everpan/mdmg/utils"
	"github.com/everpan/mdmg/web/icode"
	"github.com/everpan/mdmg/web/icode/v8runtime"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

var logger *zap.Logger

func main() {
	icode.InitLogger()
	defer logger.Sync()
	defer v8runtime.DisposeCtxPool()

	app := fiber.New()
	apiRouter := app.Group("/api")
	apiRouter.Group("/v1/icode/:module/:version?/:jsFile?/*", func(c *fiber.Ctx) error {
		fmt.Printf("%v | %v | %v | %v\n", c.Params("module"), c.Params("version"), c.Params("jsFile"), c.Queries())
		// c.Next()
		scriptFile := "web/handler/test_001.mjs"
		// zCtx := zcode.AcquireCtx(scriptFile, c)
		zCtx := v8runtime.AcquireCtx(c)
		script, e := os.ReadFile(scriptFile)
		if e != nil {
			return c.SendString(e.Error())
		}
		scriptFile = filepath.Base(scriptFile)
		r, e := zCtx.RunScript(string(script), scriptFile)
		if e != nil {
			return c.SendString(e.Error())
		}
		v, e := utils.ToGoValue(zCtx.V8Context(), r)
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
