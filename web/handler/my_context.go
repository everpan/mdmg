package handler

import (
	"github.com/everpan/mdmg/pkg/tenant"
	"github.com/everpan/mdmg/utils"
	"github.com/everpan/mdmg/v8runtime"
	"github.com/gofiber/fiber/v2"
	v8 "rogchap.com/v8go"
	"xorm.io/xorm"
)

type MyHandlerExport struct {
	Path    string
	Handler MyHandler
}

type MyHandler func(ctx *Context) error

type Context struct {
	fc            *fiber.Ctx
	tenant        *tenant.Info
	v8Ctx         *v8.Context
	db            *xorm.Engine
	ModuleVersion string
	handlers      []MyHandler
}

func (h MyHandler) WrapHandler() fiber.Handler {
	return func(fc *fiber.Ctx) error {
		ctx := AcquireContext(fc)
		return h(ctx)
	}
}

func (c *Context) RunScript(source string, origin string) (*v8.Value, error) {
	return c.v8Ctx.RunScript(source, origin)
}

func (c *Context) RunScriptRetAny(source string, origin string) (any, *v8.Value, error) {
	v, err := c.v8Ctx.RunScript(source, origin)
	if err != nil {
		return nil, nil, err
	}
	g, e := utils.ToGoValue(c.v8Ctx, v)
	return g, v, e
}

func (c *Context) CreateV8Context() *v8.Context {
	iso := v8.NewIsolate()
	icObj := v8.NewObjectTemplate(iso)
	obj := v8.NewObjectTemplate(iso)
	// ctx := AcquireCtx(fb)
	_ = obj.Set("ctx", v8runtime.ExportObject(c.fc, iso))
	_ = obj.Set("db", v8runtime.ExportXormObject(c.db, iso))
	_ = icObj.Set("__ic", obj)
	v8ctx := v8.NewContext(iso, icObj)
	// icode.logger.Info("create v8 context", zap.Any("fbCtx", fb))
	c.v8Ctx = v8ctx
	return v8ctx
}
