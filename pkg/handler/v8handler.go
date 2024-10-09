package handler

import (
	"errors"
	"fmt"
	"github.com/everpan/mdmg/pkg/config"
	"github.com/everpan/mdmg/pkg/ctx"
	"github.com/everpan/mdmg/utils"
	"github.com/gofiber/fiber/v2"
	"io/fs"
	"os"
	"path/filepath"
	v8 "rogchap.com/v8go"
	"strings"
)

var ICoderHandler = ctx.IcPathHandler{
	Path:    "/v1/icode/:modVer/:jsFile/*",
	Handler: icodeHandler,
}

func icodeHandler(c *ctx.IcContext) error {
	fc := c.FiberCtx()
	movVer := fc.Params("modVer")
	c.SetModuleVersion(movVer)
	fName := fc.Params("jsFile")
	subFile := fc.Params("*1")
	var shortFileName string
	if len(subFile) == 0 {
		shortFileName = filepath.Join(movVer, config.DefaultConfig.JSModuleBeckEndDir, fName+".js")
	} else {
		subs := strings.Split(subFile, "/")
		substr := filepath.Join(subs...)
		shortFileName = filepath.Join(movVer, config.DefaultConfig.JSModuleBeckEndDir, fName, substr+".js")
	}
	var err error
	var r1, r2, output *v8.Value
	r1, err = runScriptByFileShortName(c, shortFileName)
	if err == nil {
		defer r1.Release()
		r2, err = runMethodScript(fc.Method(), r1, c.V8Ctx())
		if err == nil {
			defer r2.Release()
			var o *v8.Object
			o, err = r2.AsObject()
			if err == nil {
				output, err = o.Get("output")
				if err == nil {
					if output.IsNullOrUndefined() {
						err = errors.New("output object is not found in response")
					} else {
						var gv any
						gv, err = utils.ToGoValue(c.V8Ctx(), output)
						if err == nil {

							resp := ctx.ICodeResponse{
								Code: 0,
								Data: gv,
							}
							fc.Response().Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
							return fc.Send(resp.Marshal())
						}
					}
				}
			}
		}
	}
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			err = errors.New(fmt.Sprintf("can not find %v", shortFileName))
			_ = ctx.SendError(fc, fiber.StatusNotFound, err)
		} else {
			return ctx.SendInternalServerError(fc, err)
		}
	}
	return nil
}

func runScriptByFileShortName(ctx *ctx.IcContext, shortFileName string) (*v8.Value, error) {
	scriptFile := filepath.Join(config.DefaultConfig.JSModuleRootPath, shortFileName)
	scriptContext, err := os.ReadFile(scriptFile)
	if err != nil {
		return nil, err
	}
	return ctx.RunScript(string(scriptContext), shortFileName)
}

func runMethodScript(method string, script *v8.Value, ctx *v8.Context) (*v8.Value, error) {
	m := strings.ToLower(method)
	if m == "delete" {
		m = "del"
	}
	scriptObj, e := script.AsObject()
	if e != nil {
		return nil, e
	}
	methodVal, e := scriptObj.Get(m)
	if e != nil {
		return nil, e
	}
	methodFun, e := methodVal.AsFunction()
	if e != nil {
		return nil, errors.New(fmt.Sprintf("not found the handler of method(%v), %v", method, e.Error()))
	}
	return methodFun.Call(ctx.Global())
}
