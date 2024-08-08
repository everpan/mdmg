package handler

import (
	"encoding/json"
	"github.com/everpan/mdmg/utils"
	"github.com/everpan/mdmg/v8runtime"
	"github.com/everpan/mdmg/web/config"
	"github.com/gofiber/fiber/v2"
	"os"
	"path/filepath"
	v8 "rogchap.com/v8go"
	"strings"
)

var ICoderHandler = PathHandler{
	Path: "/v1/icode/:modVer/:jsFile/*",
	Handler: func(fc *fiber.Ctx) error {
		zCtx := v8runtime.AcquireCtx(fc)
		zCtx.ModuleVersion = fc.Params("modVer")
		fName := fc.Params("jsFile")
		scriptFile := filepath.Join(config.DefaultConfig.JSModuleRootPath, zCtx.ModuleVersion, fName+".js")
		var err error
		var r1, r2 *v8.Value
		r1, err = runFileScript(zCtx, scriptFile)
		if err == nil {
			defer r1.Release()
			r2, err = runMethodScript(fc.Method(), r1, zCtx.V8Ctx())
			if err == nil {
				defer r2.Release()
				var gv any
				gv, err = utils.ToGoValue(zCtx.V8Context(), r2)
				if err == nil {
					var jv []byte
					jv, err = json.Marshal(gv)
					if err != nil {
						return fc.Send(jv)
					}
				}
			}
		}
		if err != nil {
			return SendInternalServerError(fc, err)
		}
		return nil
	},
}

func runMethodScript(method string, script *v8.Value, ctx *v8.Context) (*v8.Value, error) {
	m := strings.ToLower(method)
	if method == "delete" {
		method = "del"
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
		return nil, e
	}
	return methodFun.Call(ctx.Global())
}

func runFileScript(zCtx *v8runtime.Ctx, scriptFile string) (*v8.Value, error) {
	script, err := os.ReadFile(scriptFile)
	if err != nil {
		return nil, err
	}
	scriptFile = filepath.Base(scriptFile)
	return zCtx.RunScript(string(script), scriptFile)
}
