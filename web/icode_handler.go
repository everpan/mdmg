package main

import (
	"encoding/json"
	"github.com/everpan/mdmg/utils"
	"github.com/everpan/mdmg/web/icode/v8runtime"
	"github.com/gofiber/fiber/v2"
	"os"
	"path/filepath"
)

// /v1/icode/:modVer/:jsFile/*
func icodeHandler(c *fiber.Ctx) error {
	// scriptFile := "web/handler/test_001.mjs"
	// zCtx := zcode.AcquireCtx(scriptFile, c)
	root := "web/handler"
	zCtx := v8runtime.AcquireCtx(c)
	modVer := c.Params("modVer")
	fName := c.Params("jsFile")
	scriptFile := filepath.Join(root, modVer, fName+".js")

	zCtx.Module, zCtx.Version = utils.SplitModuleVersion(c.Params("modVer"))
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
	jv, e := json.Marshal(v)
	if e != nil {
		return c.SendString(e.Error())
	}
	return c.Send(jv)
}
