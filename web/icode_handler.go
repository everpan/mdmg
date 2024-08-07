package main

import (
	"encoding/json"
	"github.com/everpan/mdmg/utils"
	"github.com/everpan/mdmg/web/icode/v8runtime"
	"github.com/gofiber/fiber/v2"
	"os"
	"path/filepath"
	"strings"
)

// /v1/icode/:modVer/:jsFile/*
func icodeHandler(fc *fiber.Ctx) error {
	root := "web/handler"
	zCtx := v8runtime.AcquireCtx(fc)
	modVer := fc.Params("modVer")
	fName := fc.Params("jsFile")
	scriptFile := filepath.Join(root, modVer, fName+".js")

	zCtx.Module, zCtx.Version = utils.SplitModuleVersion(fc.Params("modVer"))
	script, e := os.ReadFile(scriptFile)
	if e != nil {
		return SendInternalServerError(fc, e)
	}
	scriptFile = filepath.Base(scriptFile)
	r, e := zCtx.RunScript(string(script), scriptFile)
	if e != nil {
		return SendInternalServerError(fc, e)
	}
	scriptObj, e := r.AsObject()
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
	r, e = methodFun.Call(zCtx.V8Context().Global())
	if e != nil {
		return SendInternalServerError(fc, e)
	}
	v, e := utils.ToGoValue(zCtx.V8Context(), r)
	if e != nil {
		return SendInternalServerError(fc, e)
	}
	jv, e := json.Marshal(v)
	if e != nil {
		return SendInternalServerError(fc, e)
	}
	return fc.Send(jv)
}

func SendInternalServerError(fc *fiber.Ctx, err error) error {
	return SendError(fc, fiber.StatusInternalServerError, err)
}
func SendError(fc *fiber.Ctx, status int, e error) error {
	fc.SendStatus(status)
	resp := NewICodeResponse(-1, e.Error(), nil)
	return fc.Send(resp.Marshal())
}
