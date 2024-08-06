package icode

import (
	"github.com/everpan/mdmg/web/icode/v8runtime"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	v8 "rogchap.com/v8go"
	"xorm.io/xorm"
)

type Ctx struct {
	fbCtx *fiber.Ctx
	v8Ctx *v8.Context
	orm   *xorm.Engine
}

// 每个`fiber.Ctx`对应一个
var pool = map[*fiber.Ctx]*Ctx{}

func CreateV8Ctx(fb *fiber.Ctx) *v8.Context {
	iso := v8.NewIsolate()
	obj := v8.NewObjectTemplate(iso)
	_ = obj.Set("icode", v8runtime.ExportObject(fb, iso))

	v8ctx := v8.NewContext(iso, obj)
	logger.Info("create v8 context", zap.Any("fbCtx", fb))
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
	//logger.Info("get ctx",
	//	zap.Uintptr("fb", uintptr(unsafe.Pointer(fb))),
	//	zap.Uintptr("ctx", uintptr(unsafe.Pointer(ctx))))
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
	fc.orm = nil
}

func (fc *Ctx) SetFiberCtx(ctx *fiber.Ctx) {
	fc.fbCtx = ctx
}

func (fc *Ctx) GetFiberCtx() *fiber.Ctx {
	return fc.fbCtx
}

func (fc *Ctx) RunScript(source string, origin string) (*v8.Value, error) {
	return fc.v8Ctx.RunScript(source, origin)
}

func (fc *Ctx) V8Context() *v8.Context {
	return fc.v8Ctx
}

func DisposeCtxPool() {
	for _, ctx := range pool {
		ctx.Dispose()
	}
}
