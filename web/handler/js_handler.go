package handler

import (
	"encoding/json"
	"github.com/everpan/mdmg/utils"
	"github.com/everpan/mdmg/v8runtime"
	"github.com/everpan/mdmg/web/config"
	"github.com/gofiber/fiber/v2"
	"os"
	"path/filepath"
	"strings"
)

var ICoderHandler = PathHandler{
	Path: "/v1/icode/:modVer/:jsFile/*",
	Handler: func(fc *fiber.Ctx) error {
		zCtx := v8runtime.AcquireCtx(fc)
		zCtx.ModuleVersion = fc.Params("modVer")
		fName := fc.Params("jsFile")
		scriptFile := filepath.Join(config.DefaultConfig.JSModuleRootPath, zCtx.ModuleVersion, fName+".js")

		// zCtx.Module, zCtx.Version = utils.SplitModuleVersion(fc.Params("modVer"))
		script, e := os.ReadFile(scriptFile)
		if e != nil {
			return SendInternalServerError(fc, e)
		}
		scriptFile = filepath.Base(scriptFile)
		r1, e := zCtx.RunScript(string(script), scriptFile)
		defer r1.Release()
		if e != nil {
			return SendInternalServerError(fc, e)
		}
		scriptObj, e := r1.AsObject()
		if e != nil {
			return SendInternalServerError(fc, e)
		}
		method := strings.ToLower(fc.Method())
		if method == "delete" {
			method = "del"
		}
		methodVal, e := scriptObj.Get(method)
		if e != nil {
			return SendInternalServerError(fc, e)
		}
		methodFun, e := methodVal.AsFunction()
		if e != nil {
			return SendInternalServerError(fc, e)
		}
		r2, e := methodFun.Call(zCtx.V8Context().Global())
		defer r2.Release()
		if e != nil {
			return SendInternalServerError(fc, e)
		}
		v, e := utils.ToGoValue(zCtx.V8Context(), r2)
		if e != nil {
			return SendInternalServerError(fc, e)
		}
		jv, e := json.Marshal(v)
		if e != nil {
			return SendInternalServerError(fc, e)
		}
		return fc.Send(jv)
	},
}
