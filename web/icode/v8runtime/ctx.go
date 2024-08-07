package v8runtime

import (
	"github.com/everpan/mdmg/utils"
	"github.com/gofiber/fiber/v2"
	v8 "rogchap.com/v8go"
	"xorm.io/xorm"
)

type Ctx struct {
	fbCtx *fiber.Ctx
	v8Ctx *v8.Context
	db    *xorm.Engine
}

// 每个`fiber.Ctx`对应一个
var pool = map[*fiber.Ctx]*Ctx{}

func CreateV8Ctx(ctx *Ctx) *v8.Context {
	iso := v8.NewIsolate()
	obj := v8.NewObjectTemplate(iso)
	// ctx := AcquireCtx(fb)
	_ = obj.Set("icode", ExportObject(ctx.fbCtx, iso))
	_ = obj.Set("db", ExportXormObject(ctx, iso))
	v8ctx := v8.NewContext(iso, obj)
	// icode.logger.Info("create v8 context", zap.Any("fbCtx", fb))
	return v8ctx
}

func createCtx(fb *fiber.Ctx) *Ctx {
	ctx := &Ctx{
		fbCtx: fb,
	}
	pool[fb] = ctx
	if ctx.v8Ctx == nil {
		ctx.v8Ctx = CreateV8Ctx(ctx)
	}
	return ctx
}

func AcquireCtx(fb *fiber.Ctx) *Ctx {
	ctx, ok := pool[fb]
	if !ok {
		ctx = createCtx(fb)
	}
	return ctx
}

func (c *Ctx) Dispose() {
	if c.v8Ctx != nil {
		iso := c.v8Ctx.Isolate()
		c.v8Ctx.Close()
		iso.Dispose()
	}
	c.fbCtx = nil
	c.v8Ctx = nil
	c.db = nil
}

func (c *Ctx) FiberCtx() *fiber.Ctx {
	return c.fbCtx
}

func (c *Ctx) Engine() *xorm.Engine {
	return c.db
}

func (c *Ctx) V8Ctx() *v8.Context {
	return c.v8Ctx
}

func (c *Ctx) RunScript(source string, origin string) (*v8.Value, error) {
	return c.v8Ctx.RunScript(source, origin)
}

func (c *Ctx) RunScriptRetAny(source string, origin string) (any, *v8.Value, error) {
	v, err := c.v8Ctx.RunScript(source, origin)
	if err != nil {
		return nil, nil, err
	}
	g, e := utils.ToGoValue(c.v8Ctx, v)
	return g, v, e
}

func (c *Ctx) V8Context() *v8.Context {
	return c.v8Ctx
}

func DisposeCtxPool() {
	for _, ctx := range pool {
		ctx.Dispose()
	}
}
