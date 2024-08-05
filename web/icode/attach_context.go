package icode

import (
	"errors"
	"fmt"
	"github.com/everpan/mdmg/utils"
	"github.com/gofiber/fiber/v2"
	v8 "rogchap.com/v8go"
)

func AttachContext(ftx *fiber.Ctx, jxt *v8.Context) error {
	//queryStrTpl := v8.FunctionTemplate{}
	ftx.Get("a")
	// off := unsafe.Offsetof(fiber.Get)

	fnTmpl := v8.NewFunctionTemplate(jxt.Isolate(), func(info *v8.FunctionCallbackInfo) *v8.Value {
		args, err := utils.ToGoValues(jxt, info.Args())
		if err != nil {
			return nil
		}
		// len(args)
		fmt.Println(args)
		// this = info.This() // return this fn object
		return nil
	})
	if fnTmpl == nil {
		return errors.New("expected FunctionTemplate, but got <nil>")
	}
	_ = fnTmpl.GetFunction(jxt)

	// objTmpl := v8.NewObjectTemplate(jxt.Isolate())
	// objTmpl.

	jxt.Global()
	return nil
}
