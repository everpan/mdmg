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

func CreateV8Ctx(fb *fiber.Ctx) *v8.Context {
	iso := v8.NewIsolate()
	obj := v8.NewObjectTemplate(iso)
	ctx := AcquireCtx(fb)
	_ = obj.Set("icode", ExportObject(fb, iso))
	_ = obj.Set("db", ExportXormObject(ctx, iso))
	v8ctx := v8.NewContext(iso, obj)
	// icode.logger.Info("create v8 context", zap.Any("fbCtx", fb))
	return v8ctx
}

func createCtx(fb *fiber.Ctx) *Ctx {
	return &Ctx{
		fbCtx: fb,
		v8Ctx: CreateV8Ctx(fb),
	}
}

func AcquireCtx(fb *fiber.Ctx) *Ctx {
	ctx, ok := pool[fb]
	if !ok {
		ctx = createCtx(fb)
		pool[fb] = ctx
	}
	//logger.Info("get _ctx",
	//	zap.Uintptr("fb", uintptr(unsafe.Pointer(fb))),
	//	zap.Uintptr("_ctx", uintptr(unsafe.Pointer(_ctx))))
	return ctx
}

func (fc *Ctx) Dispose() {
	if fc.v8Ctx != nil {
		iso := fc.v8Ctx.Isolate()
		fc.v8Ctx.Close()
		iso.Dispose()
	}
	fc.fbCtx = nil
	fc.v8Ctx = nil
	fc.db = nil
}

func (fc *Ctx) FiberCtx() *fiber.Ctx {
	return fc.fbCtx
}

func (fc *Ctx) Engine() *xorm.Engine {
	return fc.db
}

func (fc *Ctx) V8Ctx() *v8.Context {
	return fc.v8Ctx
}

func (fc *Ctx) RunScript(source string, origin string) (*v8.Value, error) {
	return fc.v8Ctx.RunScript(source, origin)
}

func (fc *Ctx) RunScriptRetAny(source string, origin string) (any, *v8.Value, error) {
	v, err := fc.v8Ctx.RunScript(source, origin)
	if err != nil {
		return nil, nil, err
	}
	g, e := utils.ToGoValue(fc.v8Ctx, v)
	return g, v, e
}

func (fc *Ctx) V8Context() *v8.Context {
	return fc.v8Ctx
}

func DisposeCtxPool() {
	for _, ctx := range pool {
		ctx.Dispose()
	}
}
