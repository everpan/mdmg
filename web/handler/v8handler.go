package handler

import (
	"errors"
	"fmt"
	"github.com/everpan/mdmg/utils"
	"github.com/everpan/mdmg/v8runtime"
	"github.com/everpan/mdmg/web/config"
	"github.com/gofiber/fiber/v2"
	"os"
	"path/filepath"
	v8 "rogchap.com/v8go"
	"strings"
)

func runScriptByFileShortName(ctx *v8runtime.Ctx, shortFileName string) (*v8.Value, error) {
	scriptFile := filepath.Join(config.DefaultConfig.JSModuleRootPath, shortFileName)
	scriptContext, err := os.ReadFile(scriptFile)
	if err != nil {
		return nil, err
	}
	return ctx.RunScript(string(scriptContext), shortFileName)
}

func icodeHandler(fc *fiber.Ctx) error {
	zCtx := v8runtime.AcquireCtx(fc)
	zCtx.ModuleVersion = fc.Params("modVer")
	fName := fc.Params("jsFile")
	shortFileName := filepath.Join(zCtx.ModuleVersion, fName+".js")
	var err error
	var r1, r2, output *v8.Value
	r1, err = runScriptByFileShortName(zCtx, shortFileName)
	if err == nil {
		defer r1.Release()
		r2, err = runMethodScript(fc.Method(), r1, zCtx.V8Ctx())
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
						gv, err = utils.ToGoValue(zCtx.V8Ctx(), output)
						if err == nil {
							resp := ICodeResponse{
								Code: 0,
								Data: gv,
							}
							return fc.Send(resp.Marshal())
						}
					}
				}
			}
		}
	}
	if err != nil {
		return SendInternalServerError(fc, err)
	}
	return nil
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
